package repositories

import (
	"context"
	"fmt"

	"github.com/AliTr404/T-MO/internal/dto"
	"github.com/AliTr404/T-MO/internal/models"
	"github.com/AliTr404/T-MO/pkg/derrors"
	"github.com/AliTr404/T-MO/pkg/dtime"
	"github.com/AliTr404/T-MO/pkg/tol"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	UserRepository interface {
		Insert(user *models.User) (primitive.ObjectID, error)
		UpdateOne(filter interface{}, update interface{}) error
		FindOneAndUpdate(filter interface{}, update interface{}) (*models.User, error)
		FindByUserID(userID string) (*dto.User, error)
		FindByUsername(username string) (*models.User, error)
		FindByEmail(email string) (*models.User, error)
		FindByEmailOrUsername(find string) (*models.User, error)
	}
	userRepository struct {
		database *mongo.Client
	}
)

func NewUserRepository(db *mongo.Client) UserRepository {
	return &userRepository{
		database: db,
	}
}

// Insert the user data to  the database
func (db *userRepository) Insert(user *models.User) (primitive.ObjectID, error) {
	collection := db.database.Database("TMO").Collection("users")
	user.CreatedAt = dtime.Now()
	user.UpdatedAt = dtime.Now()
	result, err := collection.InsertOne(context.Background(), &user)
	if err != nil {
		return result.InsertedID.(primitive.ObjectID), err
	}
	tol.TMessage(fmt.Sprintf("Repository (User) => Insert: %v", result.InsertedID))
	return result.InsertedID.(primitive.ObjectID), nil
}

// UpdateOne Without return any user model
func (db *userRepository) UpdateOne(filter interface{}, update interface{}) error {
	collection := db.database.Database("TMO").Collection("users")
	_, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return derrors.New(derrors.KindInvalid, "در بروزرسانی اطلاعات به مشکل برخوردیم")
	}
	tol.TMessage(fmt.Sprintf("Repository (User) => UpdateOne: %v", filter))
	return nil
}

// FindOneAndUpadate return the user model
func (db *userRepository) FindOneAndUpdate(filter interface{}, update interface{}) (*models.User, error) {
	var user *models.User

	collection := db.database.Database("TMO").Collection("users")
	result := collection.FindOneAndUpdate(context.Background(), filter, update).Decode(&user)
	if result != nil {
		return nil, derrors.New(derrors.KindNotFound, "کاربر مورد نظر پیدا نشد")
		// return nil, errors.New("کاربر مورد نظر پیدا نشد")
	}
	tol.TMessage(fmt.Sprintf("Repository (User) => FindOneAndUpdate: %v", filter))

	return user, nil
}

// FindByUserID return the user model
func (db *userRepository) FindByUserID(userID string) (*dto.User, error) {
	var user *dto.User
	collection := db.database.Database("TMO").Collection("users")
	userid, _ := primitive.ObjectIDFromHex(userID)
	res := collection.FindOne(context.Background(), bson.M{"_id": userid}).Decode(&user)
	if res != nil {
		return nil, derrors.New(derrors.KindNotFound, "کاربر مورد نظر پیدا نشد")
	}
	tol.TMessage(fmt.Sprintf("Repository (User) => FindByUserID: %v", userID))

	return user, nil

}

// FindByUserName return the user model
func (db *userRepository) FindByUsername(username string) (*models.User, error) {
	var user *models.User
	collection := db.database.Database("TMO").Collection("users")
	res := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if res != nil {
		return nil, derrors.New(derrors.KindNotFound, "کاربر مورد نظر پیدا نشد")
	}
	tol.TMessage(fmt.Sprintf("Repository (User) => FindByUsername: %v", user.Username))

	return user, nil

}

// FindByEmail return the user model
func (db *userRepository) FindByEmail(email string) (*models.User, error) {
	var user *models.User
	collection := db.database.Database("TMO").Collection("users")
	res := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if res != nil {
		return nil, derrors.New(derrors.KindNotFound, "کاربر مورد نظر پیدا نشد")
	}
	tol.TMessage(fmt.Sprintf("Repository (User) => FindByEmail: %v", user.Email))
	return user, nil
}

// FindByEmailOrUsername the user model
func (db *userRepository) FindByEmailOrUsername(find string) (*models.User, error) {
	var user *models.User
	collection := db.database.Database("TMO").Collection("users")
	res := collection.FindOne(context.Background(), bson.M{
		"$or": bson.A{
			bson.M{"email": find},
			bson.M{"username": find},
		},
	}).Decode(&user)

	if res != nil {
		return nil, derrors.New(derrors.KindNotFound, "کاربر مورد نظر پیدا نشد")
	}
	tol.TMessage(fmt.Sprintf("Repository (User) => FindByEmailOrUsername: %v", find))
	return user, nil
}
