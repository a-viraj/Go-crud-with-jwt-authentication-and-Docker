package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         *string            `json:"name"`
	Email        *string            `json:"email"`
	Phone        *string            `json:"phone"`
	Password     *string            `json:"password"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refreshtoken"`
	UserType     *string            `json:"usertype"`
	UserId       string             `json:"userid"`
}
