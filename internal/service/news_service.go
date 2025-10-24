package service

import (
	"context"
	"fmt"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/repository"
)

type NewsService struct {
	repo  *repository.NewsRepository
	cache *CacheService
}

func NewNewsService(repo *repository.NewsRepository, cache *CacheService) *NewsService {
	return &NewsService{
		repo:  repo,
		cache: cache,
	}
}

func (s *NewsService) CreateNews(ctx context.Context, req *models.CreateNewsRequest) (*models.News, error) {
	news := &models.News{
		Title:    req.Title,
		Summary:  req.Summary,
		Content:  req.Content,
		NewsDate: req.NewsDate,
		URL:      req.URL,
		ImageURL: req.ImageURL,
	}

	err := s.repo.Create(news)
	if err != nil {
		return nil, err
	}

	// Invalidate list cache
	_ = s.cache.DeletePattern(ctx, "news:*")

	return news, nil
}

func (s *NewsService) GetNews(ctx context.Context, id int64) (*models.News, error) {
	cacheKey := fmt.Sprintf("news:%d", id)

	// Try to get from cache
	var news models.News
	err := s.cache.Get(ctx, cacheKey, &news)
	if err == nil {
		return &news, nil
	}

	// Get from database
	result, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = s.cache.Set(ctx, cacheKey, result, 10*time.Minute)

	return result, nil
}

func (s *NewsService) ListNews(ctx context.Context, limit, offset int) ([]models.News, error) {
	cacheKey := fmt.Sprintf("news:list:%d:%d", limit, offset)

	// Try to get from cache
	var newsList []models.News
	err := s.cache.Get(ctx, cacheKey, &newsList)
	if err == nil {
		return newsList, nil
	}

	// Get from database
	newsList, err = s.repo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = s.cache.Set(ctx, cacheKey, newsList, 5*time.Minute)

	return newsList, nil
}

func (s *NewsService) UpdateNews(ctx context.Context, id int64, news *models.News) error {
	err := s.repo.Update(id, news)
	if err != nil {
		return err
	}

	// Invalidate caches
	_ = s.cache.Delete(ctx, fmt.Sprintf("news:%d", id))
	_ = s.cache.DeletePattern(ctx, "news:*")

	return nil
}

func (s *NewsService) DeleteNews(ctx context.Context, id int64) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	// Invalidate caches
	_ = s.cache.Delete(ctx, fmt.Sprintf("news:%d", id))
	_ = s.cache.DeletePattern(ctx, "news:*")

	return nil
}
