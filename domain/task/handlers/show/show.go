package handlerShowTask

import (
	"github.com/gin-gonic/gin"
	gpc "github.com/restuwahyu13/go-playground-converter"
	showControllerTask "mbf5923.com/todo/domain/task/controllers/show"
	resourceTask "mbf5923.com/todo/domain/task/resources"
	util "mbf5923.com/todo/utils"
	"net/http"
)

type handler struct {
	service showControllerTask.Service
}

func NewHandlerShowTask(service showControllerTask.Service) *handler {
	return &handler{service: service}
}

func (h *handler) ShowTaskHandler(ctx *gin.Context) {
	var input showControllerTask.InputShow
	input.Id = ctx.Param("id")

	config := gpc.ErrorConfig{
		Options: []gpc.ErrorMetaConfig{
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Id",
				Message: "id is required on param",
			},
		},
	}

	errResponse, errCount := util.GoValidator(&input, config.Options)

	if errCount > 0 {
		util.ValidatorErrorResponse(ctx, http.StatusBadRequest, http.MethodGet, errResponse)
		return
	}

	showTask, errShowTask := h.service.ShowTaskService(ctx, &input)

	switch errShowTask {

	case "RESULT_STUDENT_NOT_FOUND_404":
		util.APIResponse(ctx, "Task data is not exist or deleted", http.StatusNotFound, http.MethodGet, nil)
		return

	default:
		response := resourceTask.ResourceTask{
			Id:          showTask.ID,
			Title:       showTask.Title,
			Description: showTask.Description,
		}
		util.APIResponse(ctx, "Show Task data successfully", http.StatusOK, http.MethodGet, response)
	}
}
