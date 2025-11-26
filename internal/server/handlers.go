package server

import (
	"encoding/json"
	"net/http"
)

func (r *Router) handleHealth(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (r *Router) handleQuestions(w http.ResponseWriter, req *http.Request) {
}

func (r *Router) handleQuestionByID(w http.ResponseWriter, req *http.Request) {
}

func (r *Router) handleAnswers(w http.ResponseWriter, req *http.Request) {
}
