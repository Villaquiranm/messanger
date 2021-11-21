package grpc

import (
	"messager/client/grpc/sender"
	"messager/client/stdin"

	pb "messager/grpc/generated"

	"google.golang.org/grpc"
)

type Client struct {
	connexion *grpc.ClientConn
	sender    *sender.MessageSender
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) InitializeConexion(addr string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	chatConnexion := pb.NewChatClient(conn)
	c.connexion = conn
	c.sender = sender.NewMessageSender(stdin.NewReader(), chatConnexion)
	go c.sender.ForwardMessages()
	return nil
}

func (c *Client) Close() error {
	return c.connexion.Close()
}
