package reading_time

type CreateReadingTimeRequest struct {
	Minutes          int     `json:"minutes" binding:"required"`
	TextLength       int     `json:"text_length" binding:"required"`
	EstimatedSeconds float64 `json:"estimated_seconds" binding:"required"`
	WordCount        int     `json:"word_count" binding:"required"`
	Type             string  `json:"type" binding:"required,oneof=Blog"`
}

type UpdateReadingTimeRequest struct {
	ID               int     `json:"id" binding:"required"`
	Minutes          int     `json:"minutes" binding:"required"`
	TextLength       int     `json:"text_length" binding:"required"`
	EstimatedSeconds float64 `json:"estimated_seconds" binding:"required"`
	WordCount        int     `json:"word_count" binding:"required"`
	Type             string  `json:"type" binding:"required,oneof=Blog"`
}

type ReadingTimeResponse struct {
	ID               int     `json:"id"`
	Minutes          int     `json:"minutes"`
	TextLength       int     `json:"text_length"`
	EstimatedSeconds float64 `json:"estimated_seconds"`
	WordCount        int     `json:"word_count"`
	Type             string  `json:"type"`
	CreatedAt        string  `json:"created_at"`
}

type ReadingTimeUpdateResponse struct {
	ID               int     `json:"id"`
	Minutes          int     `json:"minutes"`
	TextLength       int     `json:"text_length"`
	EstimatedSeconds float64 `json:"estimated_seconds"`
	WordCount        int     `json:"word_count"`
	Type             string  `json:"type"`
}

type ReadingTimeDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

func ToReadingTimeResponse(p ReadingTime) ReadingTimeResponse {
	return ReadingTimeResponse{
		ID:               p.ID,
		Minutes:          p.Minutes,
		TextLength:       p.TextLength,
		EstimatedSeconds: p.EstimatedSeconds,
		WordCount:        p.WordCount,
		Type:             p.Type,
		CreatedAt:        p.CreatedAt.Format("2006-01-02"),
	}
}

func ToReadingTimeUpdateResponse(p ReadingTime) ReadingTimeUpdateResponse {
	return ReadingTimeUpdateResponse{
		ID:               p.ID,
		Minutes:          p.Minutes,
		TextLength:       p.TextLength,
		EstimatedSeconds: p.EstimatedSeconds,
		WordCount:        p.WordCount,
		Type:             p.Type,
	}
}
