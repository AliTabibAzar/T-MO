package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Video struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Caption   string               `bson:"caption" json:"caption"`
	Path      string               `bson:"path" json:"-"`
	Thumbnail string               `bson:"thumbnail" json:"thumbnail"`
	User      primitive.ObjectID   `bson:"user,omitempty" json:"user"`
	Likes     []primitive.ObjectID `bson:"likes,omitempty" json:"likes"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}
