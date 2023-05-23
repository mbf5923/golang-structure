package showControllerTask

import (
	"github.com/gin-gonic/gin"
	modelTask "mbf5923.com/todo/domain/task/models"
	"strconv"
)

type Service interface {
	ShowTaskService(ctx *gin.Context, input *InputShow) (*modelTask.EntityTask, string)
}

type service struct {
	repository Repository
}

func NewServiceShow(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) ShowTaskService(ctx *gin.Context, input *InputShow) (*modelTask.EntityTask, string) {

	id, err := strconv.ParseUint(input.Id, 0, 32)
	errorCode := make(chan string, 1)
	if err != nil {
		errorCode <- "RESULT_STUDENT_NOT_FOUND_404"
		return nil, <-errorCode
	}

	task := modelTask.EntityTask{
		ID: uint(id),
	}

	resultShowTask, errShowTask := s.repository.ShowTaskRepository(ctx, &task)

	return resultShowTask, errShowTask
}
