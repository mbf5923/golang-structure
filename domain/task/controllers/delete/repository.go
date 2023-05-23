package deleteControllerTask

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	modelTask "mbf5923.com/todo/domain/task/models"
	modelUser "mbf5923.com/todo/domain/user/models"
)

type Repository interface {
	DeleteTaskRepository(ctx *gin.Context, input *modelTask.EntityTask) string
}

type repository struct {
	db *gorm.DB
}

func NewRepositoryDelete(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) DeleteTaskRepository(ctx *gin.Context, input *modelTask.EntityTask) string {

	var task modelTask.EntityTask
	db := r.db.Model(&task)
	errorCode := make(chan string, 1)
	user, _ := ctx.Get("user")
	uid := user.(*modelUser.EntityUsers).ID
	resultStudents := db.Debug().Select("*").Where("id = ?", input.ID).Where("user_id=?", uid).Delete(&task)

	if resultStudents.RowsAffected < 1 {
		errorCode <- "RESULT_TASK_NOT_FOUND_404"
		return <-errorCode
	} else {
		errorCode <- "nil"
	}

	return <-errorCode
}
