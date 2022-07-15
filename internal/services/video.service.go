package services

import (
	"github.com/AliTr404/T-MO/internal/dto"
	"github.com/AliTr404/T-MO/internal/models"
	"github.com/AliTr404/T-MO/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	VideoService interface {
		PlayVideo(videoID string) (*models.Video, error)
		UpdateVideo(videoID string, videoDto dto.Video) error
		LikeVideo(videoID string, userID string) error
		UnLikeVideo(videoID string, userID string) error
		SendComment(videoID string, userID string, commentDto *dto.Comment) error
		GetComments(videoID string) ([]bson.M, error)
		GetCommentsReply(videoID string) ([]bson.M, error)
	}
	videoService struct {
		videoRepository   repositories.VideoRepository
		commentRepository repositories.CommentRepository
	}
)

func NewVideoService(videoRepository repositories.VideoRepository, commentRepository repositories.CommentRepository) VideoService {
	return &videoService{
		videoRepository:   videoRepository,
		commentRepository: commentRepository,
	}
}

func (s *videoService) PlayVideo(videoID string) (*models.Video, error) {
	video, err := s.videoRepository.FindByID(videoID)
	if err != nil {
		return nil, err
	}
	return video, nil
}

func (s *videoService) UpdateVideo(videoID string, videoDto dto.Video) error {
	videoid, _ := primitive.ObjectIDFromHex(videoID)
	if err := s.videoRepository.UpdateOne(bson.M{"_id": videoid}, bson.M{"$set": bson.M{"caption": videoDto.Caption}}); err != nil {
		return err
	}
	return nil
}
func (s *videoService) LikeVideo(videoID string, userID string) error {
	videoid, _ := primitive.ObjectIDFromHex(videoID)
	userid, _ := primitive.ObjectIDFromHex(userID)
	if err := s.videoRepository.UpdateOne(bson.M{"_id": videoid}, bson.M{"$addToSet": bson.M{"likes": userid}}); err != nil {
		return err
	}
	return nil
}

func (s *videoService) UnLikeVideo(videoID string, userID string) error {
	videoid, _ := primitive.ObjectIDFromHex(videoID)
	userid, _ := primitive.ObjectIDFromHex(userID)
	if err := s.videoRepository.UpdateOne(bson.M{"_id": videoid}, bson.M{"$pull": bson.M{"likes": userid}}); err != nil {
		return err
	}
	return nil
}

func (s *videoService) SendComment(videoID string, userID string, commentDto *dto.Comment) error {
	videoid, _ := primitive.ObjectIDFromHex(videoID)
	from, _ := primitive.ObjectIDFromHex(userID)
	comment := &models.Comment{
		Video:         videoid,
		From:          from,
		To:            commentDto.To,
		ParentComment: commentDto.ParentComment,
		Comment:       commentDto.Comment,
	}
	if _, err := s.commentRepository.Insert(comment); err != nil {
		return err
	}
	return nil
}

func (s *videoService) GetComments(videoID string) ([]bson.M, error) {
	videoid, _ := primitive.ObjectIDFromHex(videoID)
	comments, err := s.commentRepository.GetComments(videoid)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *videoService) GetCommentsReply(commentID string) ([]bson.M, error) {
	commentid, _ := primitive.ObjectIDFromHex(commentID)
	comments, err := s.commentRepository.GetCommentsReply(commentid)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
