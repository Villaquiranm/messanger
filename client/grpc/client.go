package grpc

import (
	"context"
	"fmt"
	"messager/client/grpc/receiver"
	"messager/client/grpc/sender"
	"messager/client/stdin"

	pb "messager/grpc/generated"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	connexion *grpc.ClientConn
	sender    *sender.MessageSender
	receiver  *receiver.ServerReceiver
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
	var loggedUser *pb.LoggedUser
	// Iterate until user gives a valid user
	fmt.Println("To enter the chat please chose an username: ")
	for {
		var username string
		fmt.Scanln(&username)
		loggedUser, err = chatConnexion.Join(context.Background(), &pb.Options{Name: username})
		if err != nil {
			// Retrieve the status from the error
			st, ok := status.FromError(err)
			if !ok {
				return err
			}
			if st.Code() == codes.AlreadyExists {
				fmt.Println("That user already exist on the chat please chose another one:")
				continue
			}
		}
		// If the user logged sucessfully continue execution
		break
	}

	c.sender = sender.NewMessageSender(stdin.NewReader(), chatConnexion, loggedUser)
	c.receiver = receiver.NewServerReceiver(chatConnexion, loggedUser)
	go c.sender.ForwardMessages()
	go c.receiver.StartReceiving()

	return nil
}

func (c *Client) Close() error {
	err := c.sender.Logout()
	if err != nil {
		fmt.Printf("Error while disconnecting the user: %s\n", err.Error())
	}
	return c.connexion.Close()
}
