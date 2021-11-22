package main

import (
	"fmt"
	"messager/client/grpc"
	"messager/config"
	"os"
	"os/signal"
)

type service struct {
	client *grpc.Client
}

func main() {
	s := &service{}
	s.startService()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("SIGINT received...")
	s.stopService()
}

func (s *service) startService() {
	s.client = grpc.NewClient()
	err := s.client.InitializeConexion(fmt.Sprintf("localhost:%d", config.GRPCPort))
	if err != nil {
		fmt.Printf("Error while initializing gRPC client: %s", err.Error())
	}
}

func (s *service) stopService() {
	fmt.Println("Stopping server")
	err := s.client.Close()
	if err != nil {
		fmt.Printf("error while stopping server: %s\n", err.Error())
	}
}
