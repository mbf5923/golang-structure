package handlerLoginUser

import (
	"github.com/gin-gonic/gin"
	gpc "github.com/restuwahyu13/go-playground-converter"
	loginAuthControllerUser "mbf5923.com/todo/domain/user/controllers/auth/login"
	util "mbf5923.com/todo/utils"
	"net/http"
)

type handler struct {
	service loginAuthControllerUser.Service
}

func NewHandlerLogin(service loginAuthControllerUser.Service) *handler {
	return &handler{service: service}
}

func (h *handler) LoginHandler(ctx *gin.Context) {

	var input loginAuthControllerUser.InputLogin
	ctx.ShouldBindJSON(&input)

	config := gpc.ErrorConfig{
		Options: []gpc.ErrorMetaConfig{
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Email",
				Message: "email is required on body",
			},
			gpc.ErrorMetaConfig{
				Tag:     "email",
				Field:   "Email",
				Message: "email format is not valid",
			},
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Password",
				Message: "password is required on body",
			},
		},
	}

	errResponse, errCount := util.GoValidator(&input, config.Options)

	if errCount > 0 {
		util.ValidatorErrorResponse(ctx, http.StatusBadRequest, http.MethodPost, errResponse)
		return
	}

	resultLogin, errLogin := h.service.LoginService(&input)

	switch errLogin {

	case "LOGIN_NOT_FOUND_404":
		util.APIResponse(ctx, "User account is not registered", http.StatusNotFound, http.MethodPost, nil)
		return

	case "LOGIN_NOT_ACTIVE_403":
		util.APIResponse(ctx, "User account is not active", http.StatusForbidden, http.MethodPost, nil)
		return

	case "LOGIN_WRONG_PASSWORD_403":
		util.APIResponse(ctx, "Username or password is wrong", http.StatusForbidden, http.MethodPost, nil)
		return

	default:
		// generate accessToken sha256 from user id and email and random string and time
		//accessTokenData := map[string]interface{}{"id": resultLogin.ID, "email": resultLogin.Email}
		//accessToken, errToken := util.Sign(accessTokenData, "JWT_SECRET", 24*60*1)

		//if errToken != nil {
		//	defer logrus.Error(errToken.Error())
		//	util.APIResponse(ctx, "Generate accessToken failed", http.StatusBadRequest, http.MethodPost, nil)
		//	return
		//}

		util.APIResponse(ctx, "Login successfully", http.StatusOK, http.MethodPost, map[string]string{"accessToken": resultLogin.ApiKey})
	}
}
