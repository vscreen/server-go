package server

import (
	"context"
	"net"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/reflection"

	vplayer "github.com/vscreen/server-go/player"
	"google.golang.org/grpc"
)

type Server struct {
	playerInstance *vplayer.Player
	subscribers    *sync.Map
	curInfo        *atomic.Value
	publishMutex   *sync.Mutex
}

func New(player *vplayer.Player) (*Server, error) {
	curInfo := atomic.Value{}

	return &Server{
		playerInstance: player,
		subscribers:    &sync.Map{},
		curInfo:        &curInfo,
	}, nil
}

func (s *Server) startNotifierService() {
	infoChannel := s.playerInstance.InfoListener()
	log.Info("[server] created notifier service")

	for info := range infoChannel {
		infoGrpc := &Info{
			Title:     info.Title,
			Thumbnail: info.Thumbnail,
			Volume:    info.Volume,
			Position:  info.Position,
			Playing:   info.Playing,
		}

		log.WithFields(log.Fields{
			"title":     info.Title,
			"thumbnail": info.Thumbnail,
			"volume":    info.Volume,
			"position":  info.Position,
			"playing":   info.Playing,
		}).Debug("[server] publishing current info to subscribers")

		s.subscribers.Range(func(subscriberID, value interface{}) bool {
			subscriber := value.(chan *Info)
			subscriber <- infoGrpc
			return true
		})
		s.curInfo.Store(infoGrpc)
	}

	// close all subscriber channels
	s.subscribers.Range(func(key, value interface{}) bool {
		subscriber := value.(chan<- *Info)
		close(subscriber)
		return true
	})
	log.Info("[server] closed notifier service")
}

func (s *Server) ListenAndServe(addr string) error {
	// setup networking
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()
	RegisterVScreenServer(grpcServer, s)
	reflection.Register(grpcServer)

	go s.startNotifierService()
	log.Info("[server] started accepting clients")
	return grpcServer.Serve(lis)
}

// GRPC's vscreen implementation

func (s *Server) Auth(ctx context.Context, c *Credential) (*Status, error) {
	return &Status{
		Code: StatusCode_OK,
	}, nil
}

func (s *Server) Play(ctx context.Context, _ *Empty) (*Status, error) {
	log.Info("[server] received play request")
	s.playerInstance.Play()

	return &Status{
		Code: StatusCode_OK,
	}, nil
}

func (s *Server) Pause(ctx context.Context, _ *Empty) (*Status, error) {
	log.Info("[server] received pause request")
	s.playerInstance.Pause()

	return &Status{
		Code: StatusCode_OK,
	}, nil
}

func (s *Server) Stop(ctx context.Context, _ *Empty) (*Status, error) {
	log.Info("[server] received stop request")
	s.playerInstance.Stop()

	return &Status{
		Code: StatusCode_OK,
	}, nil
}

func (s *Server) Next(ctx context.Context, _ *Empty) (*Status, error) {
	log.Info("[server] received next request")
	s.playerInstance.Next()

	return &Status{
		Code: StatusCode_OK,
	}, nil
}

func (s *Server) Add(ctx context.Context, src *Source) (*Status, error) {
	log.Info("[server] received add request")
	s.playerInstance.Add(src.Url)

	return &Status{
		Code: StatusCode_OK,
	}, nil
}

func (s *Server) Seek(ctx context.Context, pos *Position) (*Status, error) {
	log.Info("[server] received seek request")
	s.playerInstance.Seek(pos.GetValue())

	return &Status{
		Code: StatusCode_OK,
	}, nil
}

func (s *Server) Subscribe(user *User, stream VScreen_SubscribeServer) error {
	var err error
	id := user.GetId()

	log.WithField("id", id).Info("[server] got a new subscriber")

	subscriberChan := make(chan *Info)
	s.subscribers.Store(id, subscriberChan)
	if s.curInfo.Load() != nil {
		stream.Send(s.curInfo.Load().(*Info))
	}

	for info := range subscriberChan {
		log.WithField("id", id).Info("[server] sending current info to subscriber")
		if err = stream.Send(info); err != nil {
			break
		}
	}

	s.subscribers.Delete(user.GetId())
	log.WithField("id", id).Info("[server] unsubscribed")
	return err
}

func (s *Server) Unsubscribe(ctx context.Context, user *User) (*Status, error) {
	id := user.GetId()

	s.subscribers.Delete(id)
	log.WithField("id", id).Info("[server] unsubscribed")
	return &Status{
		Code: StatusCode_OK,
	}, nil
}

var _ VScreenServer = &Server{}
