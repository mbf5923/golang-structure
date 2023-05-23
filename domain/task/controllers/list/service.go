package listControllerTask

import (
	"github.com/gin-gonic/gin"
	modelTask "mbf5923.com/todo/domain/task/models"
)

type Service interface {
	ListTaskService(ctx *gin.Context) (*[]modelTask.EntityTask, string)
}

type service struct {
	repository Repository
}

func NewServiceList(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) ListTaskService(ctx *gin.Context) (*[]modelTask.EntityTask, string) {

	resultListTask, errListTask := s.repository.ListTaskRepository(ctx)

	return resultListTask, errListTask
}
