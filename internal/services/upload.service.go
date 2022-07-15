package services

import (
	"github.com/AliTr404/T-MO/internal/models"
	"github.com/AliTr404/T-MO/internal/repositories"
	"github.com/AliTr404/T-MO/pkg/derrors"
	"github.com/AliTr404/T-MO/pkg/thumbnail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	UploadService interface {
		UploadVideo(*models.Video) (*models.Video, error)
		UploadProfilePicture(userID string, filePath string) error
	}
	uploadService struct {
		videoRepository repositories.VideoRepository
		userRepository  repositories.UserRepository
	}
)

func NewUploadService(userRepository repositories.UserRepository, videoRepository repositories.VideoRepository) UploadService {
	return &uploadService{
		videoRepository: videoRepository,
		userRepository:  userRepository,
	}
}

func (s *uploadService) UploadVideo(video *models.Video) (*models.Video, error) {
	thumbnail, err := thumbnail.NewThumbnail(video.Path, 300, 500).Generate()
	if err != nil {
		return nil, derrors.New(derrors.KindInvalid, err.Error())
	}
	video.Thumbnail = thumbnail
	videoID, err := s.videoRepository.Insert(video)
	if err != nil {
		return nil, err
	}
	return &models.Video{ID: videoID, Thumbnail: thumbnail}, nil
}

func (s *uploadService) UploadProfilePicture(userID string, filePath string) error {
	userid, _ := primitive.ObjectIDFromHex(userID)
	if err := s.userRepository.UpdateOne(bson.M{"_id": userid}, bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "profile_picture", Value: filePath}}}}); err != nil {
		return err
	}
	return nil
}
