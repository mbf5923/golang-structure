package handlerListTask

import (
	"github.com/gin-gonic/gin"
	listControllerTask "mbf5923.com/todo/domain/task/controllers/list"
	modelTask "mbf5923.com/todo/domain/task/models"
	resourceTask "mbf5923.com/todo/domain/task/resources"
	util "mbf5923.com/todo/utils"
	"net/http"
	"reflect"
)

type handler struct {
	service listControllerTask.Service
}

func NewHandlerListTask(service listControllerTask.Service) *handler {
	return &handler{service: service}
}

func (h *handler) ListTaskHandler(ctx *gin.Context) {
	listTask, errListTask := h.service.ListTaskService(ctx)

	switch errListTask {

	case "RESULT_TASK_NOT_FOUND_404":
		util.APIResponse(ctx, "You Have Not Any Task", http.StatusNotFound, http.MethodGet, nil)
		return

	default:
		var response = make([]resourceTask.ResourceTask, reflect.ValueOf(listTask).Elem().Len())
		for index, v := range reflect.ValueOf(listTask).Elem().Interface().([]modelTask.EntityTask) {
			response[index] = resourceTask.ResourceTask{
				Id:          v.ID,
				Title:       v.Title,
				Description: v.Description,
			}
		}
		util.APIResponse(ctx, "Show Task data successfully", http.StatusOK, http.MethodGet, response)
	}
}
