package application

import (
	"errors"
	"math"
	"strings"

	"travel-api/services/review/domain"
	"travel-api/services/review/infrastructure"
)

var (
	ErrInvalidReviewRequest = errors.New("invalid_review_request")
	ErrReviewNotAllowed     = errors.New("review_not_allowed")
)

type ReviewService struct {
	repo *infrastructure.SQLiteReviewRepo
}

func NewReviewService() *ReviewService {
	return &ReviewService{repo: infrastructure.NewSQLiteReviewRepo()}
}

func (s *ReviewService) List(productID int, language string) (domain.ListResult, error) {
	reviews, err := s.repo.ListByProduct(productID, language)
	if err != nil {
		return domain.ListResult{}, err
	}
	return domain.ListResult{Summary: summarize(reviews), Reviews: reviews}, nil
}

func (s *ReviewService) Create(userID string, productID int, req domain.CreateReviewRequest) (domain.Review, error) {
	if strings.TrimSpace(userID) == "" || productID <= 0 || req.Rating < 1 || req.Rating > 5 || strings.TrimSpace(req.Content) == "" {
		return domain.Review{}, ErrInvalidReviewRequest
	}
	allowed, err := s.repo.UserHasProductOrder(userID, productID, req.OrderID)
	if err != nil {
		return domain.Review{}, err
	}
	if !allowed {
		return domain.Review{}, ErrReviewNotAllowed
	}
	language := strings.TrimSpace(req.Language)
	if language == "" {
		language = "en"
	}
	review := domain.Review{
		UserID:    userID,
		ProductID: productID,
		OrderID:   req.OrderID,
		Rating:    roundScore(req.Rating),
		Scores: domain.Scores{
			Quality:   normalizedScore(req.Scores.Quality, req.Rating),
			Service:   normalizedScore(req.Scores.Service, req.Rating),
			Value:     normalizedScore(req.Scores.Value, req.Rating),
			Transport: normalizedScore(req.Scores.Transport, req.Rating),
			Family:    normalizedScore(req.Scores.Family, req.Rating),
		},
		Content:  strings.TrimSpace(req.Content),
		Images:   req.Images,
		Language: language,
		Verified: true,
	}
	return s.repo.Create(review)
}

func summarize(reviews []domain.Review) domain.Summary {
	if len(reviews) == 0 {
		return domain.Summary{}
	}
	var summary domain.Summary
	summary.Total = len(reviews)
	for _, review := range reviews {
		summary.AverageRating += review.Rating
		summary.Quality += review.Scores.Quality
		summary.Service += review.Scores.Service
		summary.Value += review.Scores.Value
		summary.Transport += review.Scores.Transport
		summary.Family += review.Scores.Family
	}
	count := float64(len(reviews))
	summary.AverageRating = roundScore(summary.AverageRating / count)
	summary.Quality = roundScore(summary.Quality / count)
	summary.Service = roundScore(summary.Service / count)
	summary.Value = roundScore(summary.Value / count)
	summary.Transport = roundScore(summary.Transport / count)
	summary.Family = roundScore(summary.Family / count)
	return summary
}

func normalizedScore(value float64, fallback float64) float64 {
	if value < 1 || value > 5 {
		return roundScore(fallback)
	}
	return roundScore(value)
}

func roundScore(value float64) float64 {
	return math.Round(value*10) / 10
}
