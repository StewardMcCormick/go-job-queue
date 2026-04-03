package server

import (
	"log"
	"net"

	pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	listener        net.Listener
	server          *grpc.Server
	jobQueueHandler pb.JobQueueServiceServer
}

func NewServer(addr string) (gRPCServer, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return gRPCServer{}, err
	}

	server := grpc.NewServer()

	return gRPCServer{
		listener: lis,
		server:   server,
	}, nil
}

func (s *gRPCServer) Run() error {
	pb.RegisterJobQueueServiceServer(s.server, s.jobQueueHandler)

	log.Printf("Server listening on %s", s.listener.Addr().String())
	if err := s.server.Serve(s.listener); err != nil {
		return err
	}

	return nil
}

func (s *gRPCServer) Stop() error {
	if s.server != nil {
		s.server.GracefulStop()
	}

	if err := s.listener.Close(); err != nil {
		return err
	}
	return nil
}
