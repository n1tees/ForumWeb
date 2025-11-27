package rdtio

import "ForumWeb/internal/models"

// QUESTIONS
func MapQuestionToListItem(q *models.Question) QuestionListItem {
	return QuestionListItem{
		ID:        q.ID,
		Text:      q.Text,
		CreatedAt: q.CreatedAt,
	}
}

func MapQuestionsToList(items []models.Question) []QuestionListItem {
	result := make([]QuestionListItem, 0, len(items))

	for i := range items {
		result = append(result, MapQuestionToListItem(&items[i]))
	}

	return result
}

func MapQuestionToResponse(q *models.Question) QuestionResponse {
	return QuestionResponse{
		ID:        q.ID,
		Text:      q.Text,
		CreatedAt: q.CreatedAt,
	}
}

func MapQuestionToWithAnswers(q *models.Question) QuestionWithAnswersResponse {
	response := QuestionWithAnswersResponse{
		ID:        q.ID,
		Text:      q.Text,
		CreatedAt: q.CreatedAt,
		Answers:   MapAnswersToShort(q.Answers),
	}

	return response
}

// ANSWERS
func MapAnswerToShort(a *models.Answer) AnswerShortResponse {
	return AnswerShortResponse{
		ID:        a.ID,
		UserID:    a.UserID,
		Text:      a.Text,
		CreatedAt: a.CreatedAt,
	}
}

func MapAnswersToShort(items []models.Answer) []AnswerShortResponse {
	result := make([]AnswerShortResponse, 0, len(items))

	for i := range items {
		result = append(result, MapAnswerToShort(&items[i]))
	}

	return result
}

func MapAnswerToResponse(a *models.Answer) AnswerResponse {
	return AnswerResponse{
		ID:        a.ID,
		UserID:    a.UserID,
		Text:      a.Text,
		CreatedAt: a.CreatedAt,
	}
}
