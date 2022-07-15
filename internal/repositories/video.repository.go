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
	VideoRepository interface {
		Insert(video *models.Video) (primitive.ObjectID, error)
		UpdateOne(filter interface{}, update interface{}) error
		FindByID(videoID string) (*models.Video, error)
		Search(search string, page uint64) ([]bson.M, error)
		GetLatestsVideo(page uint64) ([]bson.M, error)
	}
	videoRepository struct {
		database *mongo.Client
	}
)

func NewVideoRepository(db *mongo.Client) VideoRepository {
	return &videoRepository{
		database: db,
	}
}

func (db *videoRepository) Insert(video *models.Video) (primitive.ObjectID, error) {
	collection := db.database.Database("TMO").Collection("videos")
	video.CreatedAt = dtime.Now()
	video.UpdatedAt = dtime.Now()
	result, err := collection.InsertOne(context.Background(), video)
	if err != nil {
		return result.InsertedID.(primitive.ObjectID), derrors.New(derrors.KindUnexpected, "آپلود ویدیو با شکست مواجه شد !")
	}
	tol.TMessage(fmt.Sprintf("Repository (Video) => Insert: %v", result.InsertedID))
	return result.InsertedID.(primitive.ObjectID), nil
}

func (db *videoRepository) UpdateOne(filter interface{}, update interface{}) error {
	collection := db.database.Database("TMO").Collection("videos")
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if result.MatchedCount <= 0 {
		return derrors.New(derrors.KindNotFound, "ویدیو یافت نشد !")
	}
	if err != nil {
		return derrors.New(derrors.KindInvalid, "در بروزرسانی اطلاعات به مشکل برخوردیم")
	}
	tol.TMessage(fmt.Sprintf("Repository (Video) => UpdateOne: %v", filter))
	return nil
}

func (db *videoRepository) FindByID(videoID string) (*models.Video, error) {
	var video *models.Video
	collection := db.database.Database("TMO").Collection("videos")
	videoid, _ := primitive.ObjectIDFromHex(videoID)
	result := collection.FindOne(context.Background(), bson.M{"_id": videoid}).Decode(&video)
	if result != nil {
		return nil, derrors.New(derrors.KindNotFound, "ویدیو مورد نظر یافت نشد")
	}
	tol.TMessage(fmt.Sprintf("Repository (Video) => FindByID: %v", videoid))

	return video, nil
}

func (db *videoRepository) Search(search string, page uint64) ([]bson.M, error) {
	var results []bson.M
	collection := db.database.Database("TMO").Collection("videos")
	cur, err := collection.Aggregate(context.Background(), mongo.Pipeline{
		{primitive.E{Key: "$match", Value: bson.M{"caption": bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}}}}},
		{primitive.E{Key: "$sort", Value: bson.M{"created_at": -1}}},
		{primitive.E{Key: "$lookup", Value: bson.M{"from": "users", "localField": "user", "foreignField": "_id", "as": "user"}}},
		{primitive.E{Key: "$unwind", Value: "$user"}},
		{primitive.E{Key: "$project", Value: bson.M{"user.password": 0, "path": 0}}},
		{primitive.E{Key: "$skip", Value: (page - 1) * 20}},
		{primitive.E{Key: "$limit", Value: 20}}}, options.Aggregate())
	if cur.RemainingBatchLength() <= 0 {
		return nil, derrors.New(derrors.KindNotFound, "ویدیو یافت نشد !")
	}
	if err = cur.All(context.Background(), &results); err != nil {
		return nil, derrors.New(derrors.KindInvalid, "در دریافت اطلاعات به مشکل برخوردیم")
	}
	tol.TMessage(fmt.Sprintf("Repository (Video) => Search: %v", cur.ID()))
	return results, nil
}

func (db *videoRepository) GetLatestsVideo(page uint64) ([]bson.M, error) {
	var results []bson.M
	collection := db.database.Database("TMO").Collection("videos")
	cur, err := collection.Aggregate(context.Background(), mongo.Pipeline{
		{primitive.E{Key: "$sort", Value: bson.M{"created_at": -1}}},
		{primitive.E{Key: "$lookup", Value: bson.M{"from": "users", "localField": "user", "foreignField": "_id", "as": "user"}}},
		{primitive.E{Key: "$unwind", Value: "$user"}},
		{primitive.E{Key: "$project", Value: bson.M{"user.password": 0, "path": 0}}},
		{primitive.E{Key: "$skip", Value: (page - 1) * 20}},
		{primitive.E{Key: "$limit", Value: 20}},
	}, options.Aggregate())
	if cur.RemainingBatchLength() <= 0 {
		return nil, derrors.New(derrors.KindNotFound, "ویدیو یافت نشد !")
	}
	if err = cur.All(context.Background(), &results); err != nil {
		return nil, derrors.New(derrors.KindInvalid, "در دریافت اطلاعات به مشکل برخوردیم")
	}
	tol.TMessage(fmt.Sprintf("Repository (Video) => Search: %v", cur.ID()))
	return results, nil
}
