package config

import "time"

// GRPCPort Port used by the gRPC server
const GRPCPort = 50030

// BoltDirectory path of the bold database
const (
	BoltDirectory  = "chat.db"
	UserBucket     = "users"
	MessagesBucket = "messages"
)

// MessagesCheckPeriod the frequency to check the new messages from the server
var MessagesCheckPeriod = time.Second

// AdminUsername username of administrator
const AdminUsername = "admin"
