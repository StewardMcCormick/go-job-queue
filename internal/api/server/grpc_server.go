package server

import (
	"fmt"
	"net"

	pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port string `yaml:"port" env-default:"50051"`
}

type gRPCServer struct {
	listener        net.Listener
	server          *grpc.Server
	addr            string
	jobQueueHandler pb.JobQueueServiceServer
}

func NewServer(cfg Config, log *zap.Logger, jobQueueHandler pb.JobQueueServiceServer) (gRPCServer, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return gRPCServer{}, err
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryRecoveryInterceptor()),
		grpc.UnaryInterceptor(UnaryRequestIdInterceptor(log)),
	)

	return gRPCServer{
		listener:        lis,
		server:          server,
		addr:            addr,
		jobQueueHandler: jobQueueHandler,
	}, nil
}

func (s gRPCServer) Run() error {
	pb.RegisterJobQueueServiceServer(s.server, s.jobQueueHandler)

	if err := s.server.Serve(s.listener); err != nil {
		return err
	}

	return nil
}

func (s gRPCServer) Stop() error {
	if s.server != nil {
		s.server.GracefulStop()
		return nil
	}

	if err := s.listener.Close(); err != nil {
		return fmt.Errorf("listener close error: %w", err)
	}
	return nil
}

func (s gRPCServer) Addr() string {
	return s.addr
}
