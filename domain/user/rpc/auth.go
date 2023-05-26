package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	config "mbf5923.com/todo/configs"
	modelUser "mbf5923.com/todo/domain/user/models"
	auth_pb "mbf5923.com/todo/servicepb/authpb"
	util "mbf5923.com/todo/utils"
	"net"
)

type server struct {
	auth_pb.UnimplementedAuthServiceServer
	Dbconnection *gorm.DB
}

func (s *server) mustEmbedUnimplementedAuthServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (s *server) connectToDb() {
	//TODO implement me
	s.Dbconnection = config.Connection()
}

func main() {
	fmt.Println("GRPC Server is running...")

	lis, err := net.Listen("tcp", ":"+util.GodotEnv("GRPC_PORT"))
	if err != nil {
		defer logrus.Info("Running GRPC Failed")
		logrus.Fatal(err.Error())
	}
	// Make a gRPC server
	grpcServer := grpc.NewServer()

	auth_pb.RegisterAuthServiceServer(grpcServer, &server{
		Dbconnection: config.Connection(),
	})

	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatal("Failed to serve: %v", err)
	}

}

func (s *server) Auth(ctx context.Context, req *auth_pb.AuthRequest) (*auth_pb.AuthResponse, error) {
	fmt.Printf("Received Sum RPC: %v", req)

	token := req.GetToken()

	var user modelUser.EntityUsers
	err := s.Dbconnection.Table("entity_users").Where("api_key = ?", token).First(&user).Error
	//err := config.Connection().Table("entity_users").Where("api_key = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	} else {
		res := &auth_pb.AuthResponse{
			Id: uint32(user.ID),
		}
		return res, nil
	}
}
