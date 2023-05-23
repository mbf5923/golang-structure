package showControllerTask

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	modelTask "mbf5923.com/todo/domain/task/models"
	modelUser "mbf5923.com/todo/domain/user/models"
)

type Repository interface {
	ShowTaskRepository(ctx *gin.Context, input *modelTask.EntityTask) (*modelTask.EntityTask, string)
}
type repository struct {
	db *gorm.DB
}

func NewRepositoryShow(db *gorm.DB) *repository {
	return &repository{db: db}
}
func (r *repository) ShowTaskRepository(ctx *gin.Context, input *modelTask.EntityTask) (*modelTask.EntityTask, string) {

	var task modelTask.EntityTask
	db := r.db.Model(&task)
	errorCode := make(chan string, 1)
	user, _ := ctx.Get("user")
	uid := user.(*modelUser.EntityUsers).ID
	resultStudents := db.Debug().Select("*").Where("id = ?", input.ID).Where("user_id=?", uid).Find(&task)

	if resultStudents.RowsAffected < 1 {
		errorCode <- "RESULT_STUDENT_NOT_FOUND_404"
		return &task, <-errorCode
	} else {
		errorCode <- "nil"
	}

	return &task, <-errorCode
}
