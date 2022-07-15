package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	User struct {
		ID             primitive.ObjectID `bson:"_id,omitempty" json:"userID"`
		SID            string             `bson:"_sid" json:"SID"`
		FullName       string             `bson:"fullname" json:"fullname"`
		Username       string             `bson:"username" json:"username"`
		Email          string             `bson:"email"   json:"email"`
		ProfilePicture string             `bson:"profile_picture" json:"profile_picture"`
	}
)
