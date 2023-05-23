package loginAuthControllerUser

import modelUser "mbf5923.com/todo/domain/user/models"

type Service interface {
	LoginService(input *InputLogin) (*modelUser.EntityUsers, string)
}
type service struct {
	repository Repository
}

func NewServiceLogin(repository Repository) *service {
	return &service{repository: repository}
}
func (s *service) LoginService(input *InputLogin) (*modelUser.EntityUsers, string) {

	user := modelUser.EntityUsers{
		Email:    input.Email,
		Password: input.Password,
	}

	resultLogin, errLogin := s.repository.LoginRepository(&user)

	return resultLogin, errLogin
}
