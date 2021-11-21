package main

import (
	"fmt"
	"messager/config"
	"messager/server/grpc"
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
	fmt.Println("Starting server")
	s := grpc.NewServer(config.GRPCPort)
	go s.StartListening()
}

func stopService() {
	fmt.Println("Stopping server")
}
