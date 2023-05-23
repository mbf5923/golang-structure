package deleteControllerTask

import (
	"github.com/gin-gonic/gin"
	modelTask "mbf5923.com/todo/domain/task/models"
)

type Service interface {
	DeleteTaskService(ctx *gin.Context, input *Input) string
}

type service struct {
	repository Repository
}

func NewServiceDelete(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) DeleteTaskService(ctx *gin.Context, input *Input) string {

	task := modelTask.EntityTask{
		ID: input.ID,
	}

	errDeleteTask := s.repository.DeleteTaskRepository(ctx, &task)

	return errDeleteTask
}
