package createControllerTask

import (
	"github.com/gin-gonic/gin"
	modelTask "mbf5923.com/todo/domain/task/models"
)

type Service interface {
	CreateTaskService(ctx *gin.Context, input *InputCreate) (*modelTask.EntityTask, string)
}

type service struct {
	repository Repository
}

func NewServiceCreate(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) CreateTaskService(ctx *gin.Context, input *InputCreate) (*modelTask.EntityTask, string) {

	students := modelTask.EntityTask{
		Title:       input.Title,
		Description: input.Description,
	}

	resultCreateTask, errCreateTask := s.repository.CreateTaskRepository(ctx, &students)

	return resultCreateTask, errCreateTask
}
