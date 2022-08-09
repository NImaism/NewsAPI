package Service

import (
	"Newism/Model"
	"Newism/Repository"
)

type UserService interface {
	Login(UserName string, Password string) (bool, bool, error)
	CreateUser(User Model.User) error
	UpdateUser(UserName interface{}, User Model.UpdateModel) error
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

func (userService) UpdateUser(UserName interface{}, user Model.UpdateModel) error {
	Server := Repository.URepository()
	return Server.UpdateUser(UserName, user)
}
