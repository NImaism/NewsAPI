package Repository

import (
	database "Newism/Database"
	"Newism/Model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository interface {
	Login(UserName string, Password string) (bool, bool, error)
	CreateUser(user Model.User) error
}

type userRepository struct{}

func URepository() UserRepository { return userRepository{} }

func (userRepository) Login(UserName string, Password string) (bool, bool, error) {
	usersCollection := database.GetCl(database.Data, "users")

	var Result Model.User

	err := usersCollection.FindOne(context.TODO(), bson.D{{"UserName", UserName}, {"Password", Password}}).Decode(&Result)
	if err != nil {
		return false, false, err
	}

	if Result.IsAdmin == 1 {
		return true, true, nil
	}

	return true, false, nil
}

func (userRepository) CreateUser(user Model.User) error {
	UserCol := database.GetCl(database.Data, "users")

	_, err := UserCol.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}
