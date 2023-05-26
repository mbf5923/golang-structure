package main

import (
	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log"
	config "mbf5923.com/todo/configs"
	routeTask "mbf5923.com/todo/domain/task/routes"
	routeUser "mbf5923.com/todo/domain/user/routes"
	util "mbf5923.com/todo/utils"
)

func main() {

	router := SetupRouter()
	log.Fatal(router.Run(":" + util.GodotEnv("GO_PORT")))
}

func SetupRouter() *gin.Engine {
	db := config.Connection()
	router := gin.Default()
	if util.GodotEnv("GO_ENV") != "production" && util.GodotEnv("GO_ENV") != "test" {
		gin.SetMode(gin.DebugMode)
	} else if util.GodotEnv("GO_ENV") == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))
	router.Use(helmet.Default())
	router.Use(gzip.Gzip(gzip.BestCompression))

	routeUser.InitUserRoutes(db, router)
	routeTask.InitTaskRoutes(db, router)
	//route.InitAuthRoutes(db, router)
	//route.InitStudentRoutes(db, router)

	return router
}
