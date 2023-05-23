package registerAuthControllerUser

import modelUser "mbf5923.com/todo/domain/user/models"

type Service interface {
	RegisterService(input *InputRegister) (*modelUser.EntityUsers, string)
}

type service struct {
	repository Repository
}

func NewServiceRegister(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) RegisterService(input *InputRegister) (*modelUser.EntityUsers, string) {

	users := modelUser.EntityUsers{
		Fullname: input.Fullname,
		Email:    input.Email,
		Password: input.Password,
	}

	resultRegister, errRegister := s.repository.RegisterRepository(&users)

	return resultRegister, errRegister
}
