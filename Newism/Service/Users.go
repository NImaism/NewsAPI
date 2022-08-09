package Service

import (
	"Newism/Model"
	"Newism/Repository"
)

type UserService interface {
	Login(UserName string, Password string) (bool, bool, error)
	CreateUser(User Model.User) error
}

type userService struct {
}

func NewUService() UserService {
	return userService{}
}

func (userService) Login(UserName string, Password string) (bool, bool, error) {
	Service := Repository.URepository()
	return Service.Login(UserName, Password)
}

func (userService) CreateUser(user Model.User) error {
	Server := Repository.URepository()
	return Server.CreateUser(user)
}
