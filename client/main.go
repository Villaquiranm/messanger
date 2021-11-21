package main

import (
	"fmt"
	"messager/client/grpc"
	"messager/config"
	"os"
	"os/signal"
)

func main() {
	startService()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("SIGINT received...")
	stopService()
}

func startService() {
	c := grpc.NewClient()
	err := c.InitializeConexion(fmt.Sprintf("localhost:%d", config.GRPCPort))
	if err != nil {
		fmt.Printf("Error while initializing gRPC client: %s", err.Error())
	}
}

func stopService() {
	fmt.Println("Stopping server")
}
