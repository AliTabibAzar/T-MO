package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	SID            string             `bson:"_sid" json:"_"`
	FullName       string             `bson:"fullname" json:"fullname" validate:"required"`
	Username       string             `bson:"username" json:"username" validate:"required"`
	Email          string             `bson:"email" json:"email" validate:"required,email"`
	Password       string             `bson:"password"  json:"password" validate:"required"`
	ProfilePicture string             `bson:"profile_picture" json:"profile_picture"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

func (user *User) HashPassword() ([]byte, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return pass, nil
}

func (user *User) ValidatePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err
}
