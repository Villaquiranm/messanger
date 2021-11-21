package grpc

import (
	"fmt"
	"log"
	"messager/server/grpc/services"
	"net"

	grpcApi "messager/grpc/generated"

	"google.golang.org/grpc"
)

// Server class representing a gRPC server object
type Server struct {
	port int
}

// NewServer returns a new server object
func NewServer(port int) *Server {
	return &Server{port: port}
}

// StartListening initializes gRPC server and let it ready to conexions
func (server *Server) StartListening() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", server.port))
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	grpcApi.RegisterChatServer(s, services.NewChatService())
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
}
