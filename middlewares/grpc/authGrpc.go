package grpcAuthMiddleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	modelUser "mbf5923.com/todo/domain/user/models"
	authPb "mbf5923.com/todo/servicepb/authpb"
	util "mbf5923.com/todo/utils"
	"net/http"
	"strings"
)

type UnAuthorizedError struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Method  string `json:"method"`
	Message string `json:"message"`
}
type UserClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

func Auth() gin.HandlerFunc {

	return gin.HandlerFunc(func(ctx *gin.Context) {

		var errorResponse UnAuthorizedError

		errorResponse.Status = "Forbidden"
		errorResponse.Code = http.StatusForbidden
		errorResponse.Method = ctx.Request.Method
		errorResponse.Message = "Authorization is required for this endpoint"

		if ctx.GetHeader("Authorization") == "" {
			ctx.JSON(http.StatusForbidden, errorResponse)
			defer ctx.AbortWithStatus(http.StatusForbidden)
		}

		var token = strings.Split(ctx.GetHeader("Authorization"), " ")[1]
		userGrpcUrl := fmt.Sprintf("%s:%s", util.GodotEnv("USER_GRPC_HOST"), util.GodotEnv("USER_GRPC_PORT"))
		grpcClient, err := grpc.Dial(userGrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("could not connect to server: %v", err)
		}
		defer func(grpcClient *grpc.ClientConn) {
			err := grpcClient.Close()
			if err != nil {
				log.Printf("Error while calling sum RPC: %v", err)
			}
		}(grpcClient)
		authServer := authPb.NewAuthServiceClient(grpcClient)
		var user modelUser.EntityUsers
		err = doAuth(authServer, token, &user)

		errorResponse.Status = "Unathorize"
		errorResponse.Code = http.StatusUnauthorized
		errorResponse.Method = ctx.Request.Method
		errorResponse.Message = "accessToken invalid or expired"

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse)
			defer ctx.AbortWithStatus(http.StatusUnauthorized)
		} else {
			// global value result
			ctx.Set("user", &user)
			// return to next method if token is exists
			ctx.Next()
		}
	})
}

func doAuth(authServer authPb.AuthServiceClient, token string, user *modelUser.EntityUsers) error {
	req := &authPb.AuthRequest{
		Token: token,
	}

	res, err := authServer.Auth(context.Background(), req)
	if err != nil {
		log.Printf("Error while calling sum RPC: %v", err)
		return err
	}
	js, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error while calling sum RPC: %v", err)
		return err
	}
	err = json.Unmarshal(js, &user)

	return nil
}
