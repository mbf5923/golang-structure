package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
	config "mbf5923.com/todo/configs"
	modelUser "mbf5923.com/todo/domain/user/models"
	authPb "mbf5923.com/todo/servicepb/authpb"
	util "mbf5923.com/todo/utils"
	"net"
	"time"
)

type server struct {
	authPb.UnimplementedAuthServiceServer
	DbConnection    *gorm.DB
	RedisConnection *redis.Client
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
	//Grpc Server Config Middleware
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(ContextPropagationUnaryServerInterceptor()),
	}
	// Make a gRPC server
	grpcServer := grpc.NewServer(opts...)

	authPb.RegisterAuthServiceServer(grpcServer, &server{
		DbConnection:    config.Connection(),
		RedisConnection: config.RedisConnection(),
	})

	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatal("Failed to serve: %v", err)
	}

}

func (s *server) Auth(ctx context.Context, req *authPb.AuthRequest) (*authPb.AuthResponse, error) {
	token := req.GetToken()
	var user modelUser.EntityUsers

	err := s.DbConnection.Debug().Where("api_key = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	} else {
		js, marshalErr := json.Marshal(user)
		if marshalErr != nil {
			return nil, marshalErr
		}
		var res authPb.AuthResponse
		unmarshalErr := json.Unmarshal(js, &res)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
		//set in redis with time expired 10 seconds
		err := s.RedisConnection.Set(ctx, token, js, 10*time.Second).Err()
		if err != nil {
			logrus.Info("Error Save Redis: ", err.Error())
			return nil, err
		}

		return &res, nil
	}
}

func ContextPropagationUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			println(fmt.Sprintf("Called Form: %s", md.Get("serviceName")))
			ctx = metadata.NewOutgoingContext(ctx, md)
		}
		return handler(ctx, req)
	}
}
