package handlerDeleteTask

import (
	"github.com/gin-gonic/gin"
	gpc "github.com/restuwahyu13/go-playground-converter"
	deleteControllerTask "mbf5923.com/todo/domain/task/controllers/delete"
	util "mbf5923.com/todo/utils"
	"net/http"
	"strconv"
)

type handler struct {
	service deleteControllerTask.Service
}

func NewHandlerDeleteTask(service deleteControllerTask.Service) *handler {
	return &handler{service: service}
}

func (h *handler) DeleteTaskHandler(ctx *gin.Context) {
	var input deleteControllerTask.Input
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 32)

	if err != nil {

	}

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

	input.ID = uint(id)

	errDeleteTask := h.service.DeleteTaskService(ctx, &input)

	switch errDeleteTask {
	case "RESULT_TASK_NOT_FOUND_404":
		util.APIResponse(ctx, "Task data is not exist or deleted", http.StatusNotFound, http.MethodGet, nil)
		return
	default:
		util.APIResponse(ctx, "Delete Task data successfully", http.StatusOK, http.MethodGet, nil)
		return
	}

}
