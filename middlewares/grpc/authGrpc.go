package grpcAuthMiddleware

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
	"log"
	modelUser "mbf5923.com/todo/domain/user/models"
	auth_pb "mbf5923.com/todo/servicepb/authpb"
	"net/http"
	"strings"
)

type UnathorizatedError struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Method  string `json:"method"`
	Message string `json:"message"`
}
type UserClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

func Auth(db *gorm.DB) gin.HandlerFunc {

	return gin.HandlerFunc(func(ctx *gin.Context) {

		var errorResponse UnathorizatedError

		errorResponse.Status = "Forbidden"
		errorResponse.Code = http.StatusForbidden
		errorResponse.Method = ctx.Request.Method
		errorResponse.Message = "Authorization is required for this endpoint"

		if ctx.GetHeader("Authorization") == "" {
			ctx.JSON(http.StatusForbidden, errorResponse)
			defer ctx.AbortWithStatus(http.StatusForbidden)
		}

		var token = strings.Split(ctx.GetHeader("Authorization"), " ")[1]
		cc, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("could not connect to server: %v", err)
		}
		defer cc.Close()
		c := auth_pb.NewAuthServiceClient(cc)
		var user modelUser.EntityUsers
		err = doSum(c, token, &user)

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
			// return to next method if token is exist
			ctx.Next()
		}
	})
}

func doSum(c auth_pb.AuthServiceClient, token string, user *modelUser.EntityUsers) error {
	fmt.Println("Starting to do a sum RPC")

	req := &auth_pb.AuthRequest{
		Token: token,
	}

	res, err := c.Auth(context.Background(), req)
	if err != nil {
		log.Printf("Error while calling sum RPC: %v", err)
		return err
	}

	log.Printf("Response from server: %v", res.Id)

	user.ID = uint(res.Id)

	return nil
}
