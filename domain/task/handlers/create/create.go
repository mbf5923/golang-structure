package handlerCreateTask

import (
	"github.com/gin-gonic/gin"
	gpc "github.com/restuwahyu13/go-playground-converter"
	createControllerTask "mbf5923.com/todo/domain/task/controllers/create"
	util "mbf5923.com/todo/utils"
	"net/http"
)

type handler struct {
	service createControllerTask.Service
}

func NewHandlerCreateTask(service createControllerTask.Service) *handler {
	return &handler{service: service}
}

func (h *handler) CreateStudentHandler(ctx *gin.Context) {

	var input createControllerTask.InputCreate
	ctx.ShouldBindJSON(&input)

	config := gpc.ErrorConfig{
		Options: []gpc.ErrorMetaConfig{
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Title",
				Message: "title is required on body",
			},
			gpc.ErrorMetaConfig{
				Tag:     "max",
				Field:   "Title",
				Message: "max title is 255 character",
			},
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Description",
				Message: "description is required on body",
			},
			gpc.ErrorMetaConfig{
				Tag:     "max",
				Field:   "Description",
				Message: "max description is 500 character",
			},
		},
	}

	errResponse, errCount := util.GoValidator(&input, config.Options)

	if errCount > 0 {
		util.ValidatorErrorResponse(ctx, http.StatusBadRequest, http.MethodPost, errResponse)
		return
	}

	_, errCreateStudent := h.service.CreateTaskService(ctx, &input)

	switch errCreateStudent {

	case "CREATE_STUDENT_FAILED_403":
		util.APIResponse(ctx, "Create new task failed", http.StatusForbidden, http.MethodPost, nil)
		return

	default:
		util.APIResponse(ctx, "Create new task successfully", http.StatusCreated, http.MethodPost, nil)
	}
}
