package services

import (
	"net/url"
	"strconv"

	"github.com/AliTr404/T-MO/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	HomeService interface {
		Search(queryParams url.Values) ([]bson.M, error)
		GetLatestsVideo(queryParams url.Values) ([]bson.M, error)
	}
	homeService struct {
		videoRepository repositories.VideoRepository
	}
)

func NewHomeService(videoRepository repositories.VideoRepository) HomeService {
	return &homeService{
		videoRepository: videoRepository,
	}
}

func (s *homeService) Search(queryParams url.Values) ([]bson.M, error) {
	search := queryParams.Get("search")
	page, err := strconv.ParseUint(queryParams.Get("page"), 0, 8)
	if err != nil || page <= 0 {
		page = 1
	}
	results, err := s.videoRepository.Search(search, page)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *homeService) GetLatestsVideo(queryParams url.Values) ([]bson.M, error) {
	page, err := strconv.ParseUint(queryParams.Get("page"), 0, 8)
	if err != nil || page <= 0 {
		page = 1
	}
	results, err := s.videoRepository.GetLatestsVideo(page)
	if err != nil {
		return nil, err
	}
	return results, nil
}
