package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	modelUser "mbf5923.com/todo/domain/user/models"
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

		var user modelUser.EntityUsers

		err := db.Table("entity_users").Where("api_key = ?", token).First(&user).Error
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
