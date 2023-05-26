package listControllerTask

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	modelTask "mbf5923.com/todo/domain/task/models"
	modelUser "mbf5923.com/todo/domain/user/models"
)

type Repository interface {
	ListTaskRepository(ctx *gin.Context) (*[]modelTask.EntityTask, string)
}

type repository struct {
	db *gorm.DB
}

func NewRepositoryList(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) ListTaskRepository(ctx *gin.Context) (*[]modelTask.EntityTask, string) {
	var task []modelTask.EntityTask

	db := r.db.Model(&task)
	errorCode := make(chan string, 1)
	user, _ := ctx.Get("user")
	uid := user.(*modelUser.EntityUsers).ID
	println("uid:", uid)
	resultStudents := db.Debug().Select("*").Where("user_id=?", uid).Find(&task)

	if resultStudents.RowsAffected < 1 {
		errorCode <- "RESULT_TASK_NOT_FOUND_404"
		return &task, <-errorCode
	} else {
		errorCode <- "nil"
	}

	return &task, <-errorCode

}
