package main

import (
	"context"
	"log"
	"time"

	pb "Go-000/Week04/api"

	"google.golang.org/grpc"
)

var (
	addr = "127.0.0.1:8080"
)

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGetUserClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	user, err := c.GetUserById(ctx, &pb.GetUserByIdRequest{Id: 1})
	if err != nil {
		log.Fatalf("get user infor fail: %v", err)
	}
	log.Printf("user infor：%v\n", user)
}
