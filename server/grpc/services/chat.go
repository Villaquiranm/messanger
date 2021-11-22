package services

import (
	"context"
	"fmt"
	"messager/config"
	"messager/grpc/generated"
	pb "messager/grpc/generated"
	"messager/model"
	"messager/server/database"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ChatService Struct that implements gRPC Chat service
type ChatService struct {
	pb.UnimplementedChatServer
	dbManager *database.DBManager
}

// NewChatService default constructor for ChatService type
func NewChatService() pb.ChatServer {
	manager := database.NewDBManager()
	return &ChatService{dbManager: manager}
}

// GetMessages returns the messages for a given user including the ones from himself
func (s *ChatService) GetMessages(ctx context.Context, u *pb.LoggedUser) (*pb.Messages, error) {
	messages, err := s.dbManager.MessagesForUser(u.Name)
	if err != nil {
		return &generated.Messages{}, status.Error(codes.Internal, err.Error())
	}
	response := &generated.Messages{}
	response.Messages = make([]*pb.Message, 0)
	for _, m := range messages {
		response.Messages = append(response.Messages, &generated.Message{Name: m.User, Contents: m.Content})
	}
	return response, nil
}

// Join Register a new user on the server
func (s *ChatService) Join(ctx context.Context, opt *pb.Options) (*pb.LoggedUser, error) {
	// User exist reject join to chat
	if s.dbManager.UserExist(opt.Name) {
		return &pb.LoggedUser{}, status.Error(codes.AlreadyExists, fmt.Sprintf("User:%s, already exist in the chat", opt.Name))
	}
	err := s.dbManager.CreateUser(opt.Name)
	if err != nil {
		return &pb.LoggedUser{}, status.Error(codes.Internal, err.Error())
	}
	err = s.dbManager.StoreMessage(model.Message{User: config.AdminUsername, Content: fmt.Sprintf("User %s, has join the chat!!!", opt.Name)})
	if err != nil {
		return &pb.LoggedUser{}, status.Error(codes.Internal, err.Error())
	}
	return &pb.LoggedUser{Name: opt.Name}, nil
}

// Send Sends a new message on the server
func (s *ChatService) Send(ctx context.Context, m *pb.Message) (*pb.Empty, error) {
	err := s.dbManager.StoreMessage(model.Message{User: m.Name, Content: m.Contents})
	return &pb.Empty{}, err
}

// Disconnect removes an user from the server
func (s *ChatService) Disconnect(ctx context.Context, u *pb.LoggedUser) (*pb.Empty, error) {
	err := s.dbManager.DeleteUser(u.Name)
	if err != nil {
		return &pb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	err = s.dbManager.StoreMessage(model.Message{User: config.AdminUsername, Content: fmt.Sprintf("User %s, left the chat", u.Name)})
	return &pb.Empty{}, nil
}
