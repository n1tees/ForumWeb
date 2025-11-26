package server

import (
	"net/http"

	"gorm.io/gorm"
)

type Router struct {
	db *gorm.DB
}

func NewRouter(db *gorm.DB) *Router {
	return &Router{db: db}
}

func (r *Router) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", r.handleHealth)

	mux.HandleFunc("/questions", r.handleQuestions)
	mux.HandleFunc("/questions/", r.handleQuestionByID)
	mux.HandleFunc("/answers", r.handleAnswers)
}
