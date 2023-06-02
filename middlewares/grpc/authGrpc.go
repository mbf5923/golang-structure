package grpcAuthMiddleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	modelUser "mbf5923.com/todo/domain/user/models"
	authPb "mbf5923.com/todo/servicepb/authpb"
	util "mbf5923.com/todo/utils"
	"net/http"
	"strings"
)

func Auth(redisConnection *redis.Client) gin.HandlerFunc {

	return gin.HandlerFunc(func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader("Authorization")
		if authorizationHeader == "" {
			util.APIResponse(ctx, "Authorization is required for this endpoint", http.StatusForbidden, ctx.Request.Method, nil)
			defer ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		var token = strings.Split(authorizationHeader, " ")[1]
		userGrpcUrl := fmt.Sprintf("%s:%s", util.GodotEnv("USER_GRPC_HOST"), util.GodotEnv("USER_GRPC_PORT"))
		grpcClient, err := grpc.Dial(userGrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			util.APIResponse(ctx, "GRPC Server Error", http.StatusInternalServerError, ctx.Request.Method, nil)
			defer ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer func(grpcClient *grpc.ClientConn) {
			err := grpcClient.Close()
			if err != nil {
				util.APIResponse(ctx, "GRPC Server Error", http.StatusInternalServerError, ctx.Request.Method, nil)
				defer ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}(grpcClient)
		authServer := authPb.NewAuthServiceClient(grpcClient)
		var user modelUser.EntityUsers
		fromCache, err := doAuth(authServer, token, &user, *redisConnection)
		if err != nil {
			util.APIResponse(ctx, "GRPC Server Error", http.StatusInternalServerError, ctx.Request.Method, nil)
			defer ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		} else {
			// global value result
			ctx.Set("user", &user)
			//set header token checked from cache
			ctx.Header("Token-From-Cache", fmt.Sprintf("%t", fromCache))
			// return to next method if token is exists
			ctx.Next()
		}
	})
}

func doAuth(authServer authPb.AuthServiceClient, token string, user *modelUser.EntityUsers, redisConnection redis.Client) (bool, error) {
	ctx := context.Background()
	userFromRedis, err := redisConnection.Get(ctx, token).Result()
	if err == nil {
		unmarshalErr := json.Unmarshal([]byte(userFromRedis), &user)
		if unmarshalErr != nil {
			return true, err
		}
		return true, nil
	}
	req := &authPb.AuthRequest{
		Token: token,
	}

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("serviceName", "TaskService"))
	res, err := authServer.Auth(ctx, req)
	if err != nil {
		return false, err
	}
	js, err := json.Marshal(res)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(js, &user)

	return false, nil
}
