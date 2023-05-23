package loginAuthControllerUser

import (
	"gorm.io/gorm"
	modelUser "mbf5923.com/todo/domain/user/models"
	util "mbf5923.com/todo/utils"
	"strconv"
)

type Repository interface {
	LoginRepository(input *modelUser.EntityUsers) (*modelUser.EntityUsers, string)
}

type repository struct {
	db *gorm.DB
}

func NewRepositoryLogin(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) LoginRepository(input *modelUser.EntityUsers) (*modelUser.EntityUsers, string) {
	var users modelUser.EntityUsers
	db := r.db.Model(&users)
	errorCode := make(chan string, 1)

	users.Email = input.Email
	users.Password = input.Password
	checkUserAccount := db.Debug().Select("*").Where("email = ?", input.Email).Find(&users)

	if checkUserAccount.RowsAffected < 1 {
		errorCode <- "LOGIN_NOT_FOUND_404"
		return &users, <-errorCode
	}
	if !users.Active {
		errorCode <- "LOGIN_NOT_ACTIVE_403"
		return &users, <-errorCode
	}
	comparePassword := util.ComparePassword(users.Password, input.Password)
	if comparePassword != nil {
		errorCode <- "LOGIN_WRONG_PASSWORD_403"
		return &users, <-errorCode
	} else {
		errorCode <- "nil"
	}
	var userId = strconv.FormatUint(uint64(users.ID), 10)
	apiToken := util.GenerateToken(userId, users.Email)
	db.Debug().Where("ID", users.ID).Update("ApiKey", apiToken)
	users.ApiKey = apiToken
	return &users, <-errorCode
}
