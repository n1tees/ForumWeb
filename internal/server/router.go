package server

import (
	"net/http"

	"ForumWeb/internal/service"

	"gorm.io/gorm"
)

type Router struct {
	db              *gorm.DB
	questionService *service.QuestionService
	answerService   *service.AnswerService
}

func NewRouter(db *gorm.DB) *Router {
	return &Router{
		db:              db,
		questionService: service.NewQuestionService(db),
		answerService:   service.NewAnswerService(db),
	}
}

func (r *Router) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", r.handleHealth)

	mux.HandleFunc("/questions", r.handleQuestions)
	mux.HandleFunc("/questions/", r.handleQuestionByID)
	mux.HandleFunc("/answers", r.handleAnswers)
}
