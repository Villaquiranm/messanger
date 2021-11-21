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
}

func NewMessageSender(reader *stdin.Reader, client pb.ChatClient) *MessageSender {
	return &MessageSender{
		reader: reader,
		client: client,
	}
}

func (s *MessageSender) ForwardMessages() {
	for {
		m := s.reader.Read()
		ctx := context.Background()
		_, err := s.client.Send(ctx, &pb.Message{Contents: m})
		if err != nil {
			fmt.Printf("Error sending message: %s", err.Error())
		}
	}
}
