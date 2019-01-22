package server

import (
	"context"
	"net"

	"google.golang.org/grpc/reflection"

	"github.com/vscreen/server-go/player"
	"google.golang.org/grpc"
)

type Server struct {
	playerInstance player.Player
}

func New() (*Server, error) {
	p, err := player.New()
	if err != nil {
		return nil, err
	}

	return &Server{
		playerInstance: p,
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
	return nil, nil
}

func (s *Server) Next(ctx context.Context, _ *Empty) (*Status, error) {
	return nil, nil
}

func (s *Server) Add(ctx context.Context, src *Source) (*Status, error) {
	s.playerInstance.Add(src.Url)

	return &Status{
		Code: StatusCode_OK,
	}, nil
}

func (s *Server) Seek(ctx context.Context, pos *Position) (*Status, error) {
	return nil, nil
}

func (s *Server) Subscribe(user *User, stream VScreen_SubscribeServer) error {
	return nil
}

func (s *Server) Unsubscribe(ctx context.Context, user *User) (*Status, error) {
	return nil, nil
}

var _ VScreenServer = &Server{}
