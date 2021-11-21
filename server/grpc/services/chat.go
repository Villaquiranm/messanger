package services

import (
	"context"
	"fmt"
	pb "messager/grpc/generated"
)

// ChatService Struct that implements gRPC Chat service
type ChatService struct {
	pb.UnimplementedChatServer
}

// NewChatService default constructor for ChatService type
func NewChatService() pb.ChatServer {
	return &ChatService{}
}

func (s *ChatService) Chat(ctx context.Context, e *pb.Empty) (*pb.Message, error) {
	return nil, nil
}

func (s *ChatService) Join(ctx context.Context, opt *pb.Options) (*pb.LoggedUser, error) {
	return nil, nil
}

func (s *ChatService) Send(ctx context.Context, m *pb.Message) (*pb.Empty, error) {
	fmt.Printf("User send: %s\n", m.Contents)
	return &pb.Empty{}, nil
}
