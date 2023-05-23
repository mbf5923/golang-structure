package createControllerTask

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	modelTask "mbf5923.com/todo/domain/task/models"
	modelUser "mbf5923.com/todo/domain/user/models"
)

type Repository interface {
	CreateTaskRepository(ctx *gin.Context, input *modelTask.EntityTask) (*modelTask.EntityTask, string)
}

type repository struct {
	db *gorm.DB
}

func NewRepositoryCreateTask(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) CreateTaskRepository(ctx *gin.Context, input *modelTask.EntityTask) (*modelTask.EntityTask, string) {
	var task modelTask.EntityTask
	db := r.db.Model(&task)
	errorCode := make(chan string, 1)
	user, _ := ctx.Get("user")
	uid := user.(*modelUser.EntityUsers).ID
	task.UserId = uid
	task.Title = input.Title
	task.Description = input.Description
	addNewTask := db.Debug().Create(&task)

	db.Commit()

	if addNewTask.Error != nil {
		errorCode <- "CREATE_TASK_FAILED_403"
		return &task, <-errorCode
	} else {
		errorCode <- "nil"
	}

	return &task, <-errorCode
}
