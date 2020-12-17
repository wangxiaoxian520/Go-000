package main

import (
	"Go-000/Week04/internal/service"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	service.RegisterAPI(s)
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
