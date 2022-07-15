package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	Comment struct {
		From          primitive.ObjectID `bson:"from,omitempty" json:"-"`
		To            primitive.ObjectID `bson:"to,omitempty" json:"to"`
		ParentComment primitive.ObjectID `bson:"parent_comment,omitempty" json:"parent_comment"`
		Comment       string             `bson:"comment" json:"comment" validate:"required"`
	}
)
