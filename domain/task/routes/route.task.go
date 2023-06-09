package routeTask

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	createControllerTask "mbf5923.com/todo/domain/task/controllers/create"
	deleteControllerTask "mbf5923.com/todo/domain/task/controllers/delete"
	listControllerTask "mbf5923.com/todo/domain/task/controllers/list"
	showControllerTask "mbf5923.com/todo/domain/task/controllers/show"
	handlerCreateTask "mbf5923.com/todo/domain/task/handlers/create"
	handlerDeleteTask "mbf5923.com/todo/domain/task/handlers/delete"
	handlerListTask "mbf5923.com/todo/domain/task/handlers/list"
	handlerShowTask "mbf5923.com/todo/domain/task/handlers/show"
	grpcAuthMiddleware "mbf5923.com/todo/middlewares/grpc"
)

func InitTaskRoutes(db *gorm.DB, route *gin.Engine, redisConnection *redis.Client) {
	createRepository := createControllerTask.NewRepositoryCreateTask(db)
	createService := createControllerTask.NewServiceCreate(createRepository)
	createHandler := handlerCreateTask.NewHandlerCreateTask(createService)

	showRepository := showControllerTask.NewRepositoryShow(db)
	showService := showControllerTask.NewServiceShow(showRepository)
	showHandler := handlerShowTask.NewHandlerShowTask(showService)

	listRepository := listControllerTask.NewRepositoryList(db)
	listService := listControllerTask.NewServiceList(listRepository)
	listHandler := handlerListTask.NewHandlerListTask(listService)

	deleteRepository := deleteControllerTask.NewRepositoryDelete(db)
	deleteService := deleteControllerTask.NewServiceDelete(deleteRepository)
	deleteHandler := handlerDeleteTask.NewHandlerDeleteTask(deleteService)

	groupRoute := route.Group("api/v1").Use(grpcAuthMiddleware.Auth(redisConnection))
	groupRoute.POST("task", createHandler.CreateStudentHandler)
	groupRoute.GET("task/:id", showHandler.ShowTaskHandler)
	groupRoute.GET("task", listHandler.ListTaskHandler)
	groupRoute.DELETE("task/:id", deleteHandler.DeleteTaskHandler)
}
