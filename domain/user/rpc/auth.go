package main

import (
	"context"
	"encoding/json"
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
	DbConnection *gorm.DB
}

func (s *server) mustEmbedUnimplementedAuthServiceServer() {
	panic("implement me")
}

func (s *server) connectToDb() {
	s.DbConnection = config.Connection()
}

func main() {
	grpcPort := util.GodotEnv("USER_GRPC_PORT")
	fmt.Println("GRPC Server is running on port: ", grpcPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		defer logrus.Info("Running GRPC Failed")
		logrus.Fatal(err.Error())
	}
	// Make a gRPC server
	grpcServer := grpc.NewServer()

	auth_pb.RegisterAuthServiceServer(grpcServer, &server{
		DbConnection: config.Connection(),
	})

	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatal("Failed to serve: %v", err)
	}

}

func (s *server) Auth(ctx context.Context, req *auth_pb.AuthRequest) (*auth_pb.AuthResponse, error) {
	token := req.GetToken()

	var user modelUser.EntityUsers
	err := s.DbConnection.Where("api_key = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	} else {
		js, err := json.Marshal(user)
		if err != nil {
			return nil, err
		}
		var res auth_pb.AuthResponse
		json.Unmarshal(js, &res)
		return &res, nil
	}
}
