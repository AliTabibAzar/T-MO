package repositories

import (
	"context"
	"fmt"

	"github.com/AliTr404/T-MO/internal/models"
	"github.com/AliTr404/T-MO/pkg/derrors"
	"github.com/AliTr404/T-MO/pkg/dtime"
	"github.com/AliTr404/T-MO/pkg/tol"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	CommentRepository interface {
		Insert(comment *models.Comment) (primitive.ObjectID, error)
		GetComments(videoID primitive.ObjectID) ([]bson.M, error)
		GetCommentsReply(videoID primitive.ObjectID) ([]bson.M, error)
	}
	commentRepository struct {
		database *mongo.Client
	}
)

func NewCommentRepository(db *mongo.Client) CommentRepository {
	return &commentRepository{
		database: db,
	}
}

func (db *commentRepository) Insert(comment *models.Comment) (primitive.ObjectID, error) {
	collection := db.database.Database("TMO").Collection("comments")
	comment.CreatedAt = dtime.Now()
	comment.UpdatedAt = dtime.Now()
	result, err := collection.InsertOne(context.Background(), comment)
	if err != nil {
		return result.InsertedID.(primitive.ObjectID), derrors.New(derrors.KindUnexpected, "ثبت کامنت با شکست مواجه شد !")
	}
	tol.TMessage(fmt.Sprintf("Repository (Comment) => Insert: %v", result.InsertedID))
	return result.InsertedID.(primitive.ObjectID), nil
}

func (db *commentRepository) GetComments(videoID primitive.ObjectID) ([]bson.M, error) {
	var comments []bson.M
	collection := db.database.Database("TMO").Collection("comments")
	cur, err := collection.Aggregate(context.Background(), mongo.Pipeline{
		{primitive.E{Key: "$match", Value: bson.M{"video": videoID, "parent_comment": bson.M{"$exists": false}}}},
		{primitive.E{Key: "$lookup", Value: bson.M{"from": "users", "localField": "from", "foreignField": "_id", "as": "from"}}},
		{primitive.E{Key: "$unwind", Value: "$from"}},
		{primitive.E{Key: "$lookup", Value: bson.M{"from": "users", "localField": "to", "foreignField": "_id", "as": "to"}}},
		{primitive.E{Key: "$unwind", Value: "$to"}},
		{primitive.E{Key: "$project", Value: bson.M{"from.password": 0, "to.password": 0}}},
	}, options.Aggregate())
	if err != nil {
		return nil, derrors.New(derrors.KindInvalid, "در دریافت اطلاعات به مشکل برخوردیم")
	}
	if err = cur.All(context.Background(), &comments); err != nil {
		return nil, derrors.New(derrors.KindInvalid, "در دریافت اطلاعات به مشکل برخوردیم")
	}
	tol.TMessage(fmt.Sprintf("Repository (Comment) => GetCommentsByID: %v", videoID))
	return comments, nil
}

func (db *commentRepository) GetCommentsReply(commentID primitive.ObjectID) ([]bson.M, error) {
	var comments []bson.M
	collection := db.database.Database("TMO").Collection("comments")
	cur, err := collection.Aggregate(context.Background(), mongo.Pipeline{
		{primitive.E{Key: "$match", Value: bson.M{"parent_comment": commentID}}},
		{primitive.E{Key: "$lookup", Value: bson.M{"from": "users", "localField": "from", "foreignField": "_id", "as": "from"}}},
		{primitive.E{Key: "$unwind", Value: "$from"}},
		{primitive.E{Key: "$lookup", Value: bson.M{"from": "users", "localField": "to", "foreignField": "_id", "as": "to"}}},
		{primitive.E{Key: "$unwind", Value: "$to"}},
		{primitive.E{Key: "$project", Value: bson.M{"from.password": 0, "to.password": 0}}},
	}, options.Aggregate())
	if err != nil {
		return nil, derrors.New(derrors.KindInvalid, "در دریافت اطلاعات به مشکل برخوردیم")
	}
	if err = cur.All(context.Background(), &comments); err != nil {
		return nil, derrors.New(derrors.KindInvalid, "در دریافت اطلاعات به مشکل برخوردیم")
	}
	tol.TMessage(fmt.Sprintf("Repository (Comment) => GetCommentsReply: %v", commentID))
	return comments, nil
}
