package routeUser

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	loginAuthControllerUser "mbf5923.com/todo/domain/user/controllers/auth/login"
	registerAuthControllerUser "mbf5923.com/todo/domain/user/controllers/auth/register"
	handlerLoginUser "mbf5923.com/todo/domain/user/handlers/login"
	handlerRegisterUser "mbf5923.com/todo/domain/user/handlers/register"
)

func InitUserRoutes(db *gorm.DB, route *gin.Engine) {
	registerRepository := registerAuthControllerUser.NewRepositoryRegister(db)
	registerService := registerAuthControllerUser.NewServiceRegister(registerRepository)
	registerHandler := handlerRegisterUser.NewHandlerRegister(registerService)

	loginRepository := loginAuthControllerUser.NewRepositoryLogin(db)
	loginService := loginAuthControllerUser.NewServiceLogin(loginRepository)
	loginHandler := handlerLoginUser.NewHandlerLogin(loginService)

	groupRoute := route.Group("/api/v1")
	groupRoute.POST("/register", registerHandler.RegisterHandler)
	groupRoute.POST("/login", loginHandler.LoginHandler)
}
