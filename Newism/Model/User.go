package Model

import (
	"github.com/golang-jwt/jwt"
)

type JwtClaims struct {
	UserName string
	IsAdmin  bool
	jwt.StandardClaims
}

type LoginUserViewModel struct {
	UserName string `json:"UserName" validate:"required"`
	Password string `json:"Password" validate:"required"`
}

type User struct {
	Id          string `bson:"_id,omitempty"`
	UserName    string `bson:"UserName" json:"UserName" validate:"required"`
	Password    string `bson:"Password" json:"Password" validate:"required"`
	FirstName   string `bson:"FirstName" json:"FirstName" validate:"required"`
	LastName    string `bson:"LastName" json:"LastName" validate:"required"`
	Age         string `bson:"Age" json:"Age" validate:"required"`
	PhoneNumber string `bson:"PhoneNumber" json:"PhoneNumber" validate:"required"`
	Avatar      string `bson:"Avatar,omitempty" json:"Avatar" validate:"required"`
	IsAdmin     int    `bson:"IsAdmin"`
}

type UpdateModel struct {
	Password    string `bson:"Password,omitempty" json:"Password" `
	FirstName   string `bson:"FirstName,omitempty" json:"FirstName"`
	LastName    string `bson:"LastName,omitempty" json:"LastName"`
	Age         string `bson:"Age,omitempty" json:"Age"`
	PhoneNumber string `bson:"PhoneNumber,omitempty" json:"PhoneNumber"`
	Avatar      string `bson:"Avatar,omitempty" json:"Avatar"`
}
