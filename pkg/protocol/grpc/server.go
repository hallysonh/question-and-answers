package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"question-and-answers/pkg/api/question"
	service "question-and-answers/pkg/service/question"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type ApiServer struct {
	server *grpc.Server
}

func NewApiServer(db *gorm.DB) *ApiServer {
	// gRPC server startup options
	var opts []grpc.ServerOption

	// register service
	server := grpc.NewServer(opts...)
	question.RegisterQuestionServiceServer(
		server, service.NewQuestionServiceServer(db),
	)
	return &ApiServer{server: server}
}

func (s *ApiServer) Start(ctx context.Context, serverPort string) error {
	listen, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		return fmt.Errorf(
			"error trying to listen port: %s -> %v",
			serverPort, err,
		)
	}

	// start gRPC server
	slog.InfoContext(ctx, "starting gRPC server...")
	return s.server.Serve(listen)
}

func (s *ApiServer) Stop() error {
	slog.Warn("shutting down gRPC server...")
	s.server.GracefulStop()
	return nil
}
