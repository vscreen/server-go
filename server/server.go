package server

import (
	"context"
	"net"
	"sync"

	"google.golang.org/grpc/reflection"

	"github.com/vscreen/server-go/player"
	"google.golang.org/grpc"
)

type Server struct {
	playerInstance player.Player
	subscribers    *sync.Map
	curInfo        *Info
}

func New() (*Server, error) {
	p, err := player.New()
	if err != nil {
		return nil, err
	}

	return &Server{
		playerInstance: p,
		subscribers:    &sync.Map{},
		curInfo:        nil,
	}, nil
}

func (s *Server) ListenAndServe(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()
	RegisterVScreenServer(grpcServer, s)
	reflection.Register(grpcServer)

	go s.playerInstance.Start()
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
	s.subscribers.Store(user.GetId(), stream)
	infoChannel := s.playerInstance.InfoListener()

	if s.curInfo != nil {
		stream.Send(s.curInfo)
	}

	for info := range infoChannel {
		state := Info_STOPPED
		switch info.State {
		case "playing":
			state = Info_PLAYING
		case "paused":
			state = Info_PAUSED
		}

		infoGrpc := &Info{
			Title:        info.Title,
			ThumbnailURL: info.Thumbnail,
			Volume:       info.Volume,
			Position:     info.Position,
			State:        state,
		}

		s.subscribers.Range(func(key, value interface{}) bool {
			subscriber := value.(VScreen_SubscribeServer)
			if err := subscriber.Send(infoGrpc); err != nil {
				s.subscribers.Delete(key)
			}
			return true
		})
		s.curInfo = infoGrpc

		if _, ok := s.subscribers.Load(user.GetId()); !ok {
			break
		}
	}

	return nil
}

func (s *Server) Unsubscribe(ctx context.Context, user *User) (*Status, error) {
	s.subscribers.Delete(user.GetId())
	return &Status{
		Code: StatusCode_OK,
	}, nil
}

var _ VScreenServer = &Server{}
