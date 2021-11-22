package sender

import (
	"context"
	"fmt"
	"messager/client/stdin"
	pb "messager/grpc/generated"
)

type MessageSender struct {
	reader *stdin.Reader
	client pb.ChatClient
	user   *pb.LoggedUser
}

// NewMessageSender sends messages everytime the logged user send an input using the stdin
func NewMessageSender(reader *stdin.Reader, client pb.ChatClient, loggedUser *pb.LoggedUser) *MessageSender {
	return &MessageSender{
		reader: reader,
		client: client,
		user:   loggedUser,
	}
}

// ForwardMessages takes all the messages from the stdin and forwards them to the gRPC server
func (s *MessageSender) ForwardMessages() {
	fmt.Println("You are now connected, you can start chatting")
	for {
		m := s.reader.Read()
		ctx := context.Background()
		_, err := s.client.Send(ctx, &pb.Message{Contents: m, Name: s.user.Name})
		if err != nil {
			fmt.Printf("Error sending message: %s", err.Error())
		}
	}
}

// Logout disconnects the user from the chat services
func (s *MessageSender) Logout() error {
	ctx := context.Background()
	_, err := s.client.Disconnect(ctx, s.user)
	return err
}
