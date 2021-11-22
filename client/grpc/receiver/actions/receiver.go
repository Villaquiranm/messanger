package actions

import pb "messager/grpc/generated"

// Receiver receives a message and do an action with it
type Receiver interface {
	Receive([]*pb.Message, string)
}
