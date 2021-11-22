package actions

import (
	"fmt"
	"messager/config"
	pb "messager/grpc/generated"

	"github.com/fatih/color"
)

const messageTemplate = "%s: %s\n"

// ReceiverPrint type that implements receiver interface
type ReceiverPrint struct{}

// NewReceiverPrint return a new ReceiverPrint struct
func NewReceiverPrint() Receiver {
	return &ReceiverPrint{}
}

// Receive Takes a message and print it to Stdout
func (r *ReceiverPrint) Receive(messages []*pb.Message, client string) {
	for _, m := range messages {
		if m.Name == config.AdminUsername {
			c := color.New(color.FgCyan).Add(color.Underline)
			c.Print(fmt.Sprintf(messageTemplate, m.Name, m.Contents))
		} else if m.Name != client {
			fmt.Printf(messageTemplate, m.Name, m.Contents)
		}
	}
}
