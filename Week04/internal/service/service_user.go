package service

import (
	pb "Go-000/Week04/api"
	"Go-000/Week04/internal/data"
	"Go-000/Week04/internal/pkg"
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedGetUserServer
}

//注册api服务
func RegisterAPI(s *grpc.Server) {
	pb.RegisterGetUserServer(s, &Server{})
}

//api
func (s *Server) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	id := req.GetId()
	fmt.Println(id)
	//对象初始化，最好放在初始化的方法里么？
	//grpc
	biz, err := pkg.InitializeUser()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	user, err := biz.GetUserById(int(id))
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return nil, errors.Unwrap(err)
		}
		fmt.Println(err)
		return nil, err
	}
	reply := &pb.GetUserByIdResponse{
		Name: user.Name,
		Age:  int32(user.Age),
		Addr: user.Addr,
	}
	return reply, nil
}
