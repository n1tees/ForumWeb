package server

import (
	"ForumWeb/internal/rdtio"
	"ForumWeb/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
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
	switch req.Method {
	case http.MethodPost:
		r.handleCreateQuestion(w, req)
	case http.MethodGet:
		r.handleListQuestions(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleQuestionByID(w http.ResponseWriter, req *http.Request) {
	// paths:
	// - /questions/123
	// - /questions/123/answers

	path := strings.Trim(req.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	// Если это answers: /questions/{id}/answers
	if len(parts) == 3 && parts[2] == "answers" {
		if req.Method == http.MethodPost {
			r.handleCreateAnswer(w, req)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ---- обрабатываем /questions/{id} ----

	switch req.Method {
	case http.MethodGet:
		r.handleGetQuestion(w, req)
	case http.MethodDelete:
		r.handleDeleteQuestion(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleAnswers(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.handleGetAnswer(w, req)
	case http.MethodDelete:
		r.handleDeleteAnswer(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleCreateQuestion(w http.ResponseWriter, req *http.Request) {
	var body rdtio.CreateQuestionRequest

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(body.Text) == "" {
		http.Error(w, "text is required", http.StatusBadRequest)
		return
	}

	q, err := r.questionService.CreateQuestion(req.Context(), body.Text)
	if err != nil {
		if errors.Is(err, service.ErrQuestionTextEmpty) {
			http.Error(w, "text is empty", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := rdtio.MapQuestionToResponse(q)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (r *Router) handleListQuestions(w http.ResponseWriter, req *http.Request) {
	items, err := r.questionService.GetAllQuestions(req.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := rdtio.MapQuestionsToList(items)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (r *Router) handleGetQuestion(w http.ResponseWriter, req *http.Request) {
	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	q, err := r.questionService.GetQuestionByID(req.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			http.Error(w, "question not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	resp := rdtio.MapQuestionToWithAnswers(q)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (r *Router) handleDeleteQuestion(w http.ResponseWriter, req *http.Request) {
	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = r.questionService.DeleteQuestion(req.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			http.Error(w, "question not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rdtio.StatusResponse{Status: "deleted"})
}

func (r *Router) handleCreateAnswer(w http.ResponseWriter, req *http.Request) {
	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	qID, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "invalid question id", http.StatusBadRequest)
		return
	}

	var body rdtio.CreateAnswerRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(body.UserID) == "" || strings.TrimSpace(body.Text) == "" {
		http.Error(w, "user_id and text required", http.StatusBadRequest)
		return
	}

	answer, err := r.answerService.CreateAnswer(req.Context(), uint(qID), body.UserID, body.Text)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrParentNotFound):
			http.Error(w, "question not found", http.StatusNotFound)
		case errors.Is(err, service.ErrAnswerTextEmpty):
			http.Error(w, "text empty", http.StatusBadRequest)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	resp := rdtio.MapAnswerToResponse(answer)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (r *Router) handleGetAnswer(w http.ResponseWriter, req *http.Request) {
	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	aID, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "invalid answer id", http.StatusBadRequest)
		return
	}

	answer, err := r.answerService.GetAnswerByID(req.Context(), uint(aID))
	if err != nil {
		if errors.Is(err, service.ErrAnswerNotFound) {
			http.Error(w, "answer not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	resp := rdtio.MapAnswerToResponse(answer)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (r *Router) handleDeleteAnswer(w http.ResponseWriter, req *http.Request) {
	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	aID, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "invalid answer id", http.StatusBadRequest)
		return
	}

	err = r.answerService.DeleteAnswer(req.Context(), uint(aID))
	if err != nil {
		if errors.Is(err, service.ErrAnswerNotFound) {
			http.Error(w, "answer not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rdtio.StatusResponse{Status: "deleted"})
}
