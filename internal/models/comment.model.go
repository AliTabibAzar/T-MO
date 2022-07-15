package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Video         primitive.ObjectID `bson:"video,omitempty" json:"-"`
	From          primitive.ObjectID `bson:"from,omitempty" json:"from"`
	To            primitive.ObjectID `bson:"to,omitempty" json:"to"`
	ParentComment primitive.ObjectID `bson:"parent_comment,omitempty" json:"parent_comment"`
	Comment       string             `bson:"comment" json:"comment"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}
