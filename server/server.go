package server

import (
	"context"
	"net"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc/reflection"

	"github.com/vscreen/server-go/player"
	"google.golang.org/grpc"
)

type Server struct {
	playerInstance player.Player
	subscribers    *sync.Map
	curInfo        *atomic.Value
	publishMutex   *sync.Mutex
}

func New() (*Server, error) {
	curInfo := atomic.Value{}

	return &Server{
		playerInstance: nil,
		subscribers:    &sync.Map{},
		curInfo:        &curInfo,
	}, nil
}

func (s *Server) startNotifierService() {
	infoChannel := s.playerInstance.InfoListener()

	for info := range infoChannel {
		infoGrpc := &Info{
			Title:     info.Title,
			Thumbnail: info.Thumbnail,
			Volume:    info.Volume,
			Position:  info.Position,
			Playing:   info.Playing,
		}

		s.subscribers.Range(func(key, value interface{}) bool {
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
}

func (s *Server) ListenAndServe(addr string) error {
	// setup player first
	p, err := player.New()
	if err != nil {
		return err
	}
	defer p.Close()
	s.playerInstance = p
	go s.playerInstance.Start()

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
	return grpcServer.Serve(lis)
}

// GRPC's vscreen implementation

func (s *Server) Auth(ctx context.Context, c *Credential) (*Status, error) {
	return nil, nil
}

func (s *Server) Play(ctx context.Context, _ *Empty) (*Status, error) {

	s.playerInstance.Play()

	return &Status{
		Code: StatusCode_OK,
	}, nil
}

func (s *Server) Pause(ctx context.Context, _ *Empty) (*Status, error) {
	s.playerInstance.Pause()

	return &Status{
		Code: StatusCode_OK,
	}, nil
}

func (s *Server) Stop(ctx context.Context, _ *Empty) (*Status, error) {
	s.playerInstance.Stop()

	return &Status{
		Code: StatusCode_OK,
	}, nil
}

func (s *Server) Next(ctx context.Context, _ *Empty) (*Status, error) {
	s.playerInstance.Next()

	return &Status{
		Code: StatusCode_OK,
	}, nil
}

func (s *Server) Add(ctx context.Context, src *Source) (*Status, error) {
	s.playerInstance.Add(src.Url)

	return &Status{
		Code: StatusCode_OK,
	}, nil
}

func (s *Server) Seek(ctx context.Context, pos *Position) (*Status, error) {
	s.playerInstance.Seek(pos.GetValue())

	return &Status{
		Code: StatusCode_OK,
	}, nil
}

func (s *Server) Subscribe(user *User, stream VScreen_SubscribeServer) error {
	var err error

	subscriberChan := make(chan *Info)
	s.subscribers.Store(user.GetId(), subscriberChan)
	if s.curInfo.Load() != nil {
		stream.Send(s.curInfo.Load().(*Info))
	}

	for info := range subscriberChan {
		if err = stream.Send(info); err != nil {
			break
		}
	}

	s.subscribers.Delete(user.GetId())
	return err
}

func (s *Server) Unsubscribe(ctx context.Context, user *User) (*Status, error) {
	s.subscribers.Delete(user.GetId())
	return &Status{
		Code: StatusCode_OK,
	}, nil
}

var _ VScreenServer = &Server{}
