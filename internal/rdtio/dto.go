package rdtio

import "time"

// DTO QUESTIONS
// GET /questions
type QuestionListItem struct {
	ID        uint      `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

// POST /questions
type QuestionResponse struct {
	ID        uint      `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

// GET /questions/{id}
type QuestionWithAnswersResponse struct {
	ID        uint                  `json:"id"`
	Text      string                `json:"text"`
	CreatedAt time.Time             `json:"created_at"`
	Answers   []AnswerShortResponse `json:"answers"`
}

// DTO ANSWERS
// POST /questions/{id}/answers

type AnswerResponse struct {
	ID         uint      `json:"id"`
	QuestionID uint      `json:"question_id"`
	UserID     string    `json:"user_id"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
}

type AnswerShortResponse struct {
	ID        uint      `json:"id"`
	UserID    string    `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type StatusResponse struct {
	Status string `json:"status"`
}
