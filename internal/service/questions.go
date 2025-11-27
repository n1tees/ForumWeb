package service

import (
	"context"
	"errors"
	"strings"

	"ForumWeb/internal/models"

	"gorm.io/gorm"
)

type QuestionService struct {
	db *gorm.DB
}

var (
	ErrQuestionTextEmpty = errors.New("question text is empty")
	ErrQuestionNotFound  = errors.New("question not found")
)

func NewQuestionService(db *gorm.DB) *QuestionService {
	return &QuestionService{db: db}
}

func (s *QuestionService) CreateQuestion(ctx context.Context, text string) (*models.Question, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, ErrQuestionTextEmpty
	}

	q := &models.Question{
		Text: text,
	}

	if err := s.db.WithContext(ctx).Create(q).Error; err != nil {
		return nil, err
	}

	return q, nil
}

func (s *QuestionService) GetAllQuestions(ctx context.Context) ([]models.Question, error) {
	var items []models.Question

	if err := s.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (s *QuestionService) GetQuestionByID(ctx context.Context, id uint) (*models.Question, error) {
	var q models.Question

	err := s.db.WithContext(ctx).
		Preload("Answers").
		First(&q, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrQuestionNotFound
	}

	return &q, err
}

func (s *QuestionService) DeleteQuestion(ctx context.Context, id uint) error {
	res := s.db.WithContext(ctx).Delete(&models.Question{}, id)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return ErrQuestionNotFound
	}

	return nil
}
