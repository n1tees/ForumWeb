package rdtio

// DTI QUESTIONS
// POST /questions
type CreateQuestionRequest struct {
	Text string `json:"text"`
}

// DTI ANSWERS
// POST /questions/{id}/answers
type CreateAnswerRequest struct {
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}
