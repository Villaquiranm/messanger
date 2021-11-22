package receiver

import (
	"context"
	"fmt"
	"messager/client/grpc/receiver/actions"
	"messager/config"
	pb "messager/grpc/generated"
	"time"
)

// ServerReceiver consults regurlary the server and realizes and action determined by each receiver
type ServerReceiver struct {
	recivers []actions.Receiver
	client   pb.ChatClient
	user     *pb.LoggedUser
}

// NewServerReceiver return a new ServerReceiver object
func NewServerReceiver(client pb.ChatClient, user *pb.LoggedUser) *ServerReceiver {
	return &ServerReceiver{
		client:   client,
		user:     user,
		recivers: []actions.Receiver{&actions.ReceiverPrint{}},
	}
}

// StartReceiving function that periodically checks for new messages and receives them
func (s *ServerReceiver) StartReceiving() {
	for {
		time.Sleep(config.MessagesCheckPeriod)
		ctx := context.Background()
		messages, err := s.client.GetMessages(ctx, s.user)
		if err != nil {
			fmt.Printf("error while receiving messages: %s", err.Error())
			continue
		}
		if len(messages.Messages) == 0 {
			continue
		}
		for _, r := range s.recivers {
			r.Receive(messages.Messages, s.user.Name)
		}
	}
}
