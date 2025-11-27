package service

import (
	"context"
	"errors"
	"strings"

	"ForumWeb/internal/models"

	"gorm.io/gorm"
)

type AnswerService struct {
	db *gorm.DB
}

var (
	ErrAnswerTextEmpty = errors.New("answer text is empty")
	ErrAnswerNotFound  = errors.New("answer not found")
	ErrParentNotFound  = errors.New("parent question not found")
)

func NewAnswerService(db *gorm.DB) *AnswerService {
	return &AnswerService{db: db}
}

func (s *AnswerService) CreateAnswer(ctx context.Context, questionID uint, userID, text string) (*models.Answer, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, ErrAnswerTextEmpty
	}

	var q models.Question
	if err := s.db.WithContext(ctx).First(&q, questionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrParentNotFound
		}
		return nil, err
	}

	a := &models.Answer{
		QuestionID: questionID,
		UserID:     userID,
		Text:       text,
	}

	if err := s.db.WithContext(ctx).Create(a).Error; err != nil {
		return nil, err
	}

	return a, nil
}

func (s *AnswerService) GetAnswerByID(ctx context.Context, id uint) (*models.Answer, error) {
	var a models.Answer

	err := s.db.WithContext(ctx).First(&a, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAnswerNotFound
	}

	return &a, err
}

func (s *AnswerService) DeleteAnswer(ctx context.Context, id uint) error {
	res := s.db.WithContext(ctx).Delete(&models.Answer{}, id)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrAnswerNotFound
	}

	return nil
}
