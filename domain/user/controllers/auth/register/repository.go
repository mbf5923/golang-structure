package registerAuthControllerUser

import (
	"gorm.io/gorm"
	modelUser "mbf5923.com/todo/domain/user/models"
)

type Repository interface {
	RegisterRepository(input *modelUser.EntityUsers) (*modelUser.EntityUsers, string)
}
type repository struct {
	db *gorm.DB
}

func NewRepositoryRegister(db *gorm.DB) *repository {
	return &repository{db: db}
}
func (r *repository) RegisterRepository(input *modelUser.EntityUsers) (*modelUser.EntityUsers, string) {

	var users modelUser.EntityUsers
	db := r.db.Model(&users)
	errorCode := make(chan string, 1)

	checkUserAccount := db.Debug().Select("*").Where("email = ?", input.Email).Find(&users)

	if checkUserAccount.RowsAffected > 0 {
		errorCode <- "REGISTER_CONFLICT_409"
		return &users, <-errorCode
	}

	users.Fullname = input.Fullname
	users.Email = input.Email
	users.Password = input.Password

	addNewUser := db.Debug().Create(&users)
	db.Commit()

	if addNewUser.Error != nil {
		errorCode <- "REGISTER_FAILED_403"
		return &users, <-errorCode
	} else {
		errorCode <- "nil"
	}

	return &users, <-errorCode
}
