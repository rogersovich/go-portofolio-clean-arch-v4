package statistic

type CreateStatisticRequest struct {
	Likes *int   `json:"likes" binding:"required"`
	Views *int   `json:"views" binding:"required"`
	Type  string `json:"type" binding:"required,oneof=Blog Project"`
}

type UpdateStatisticRequest struct {
	Id    int    `json:"id" binding:"required"`
	Likes *int   `json:"likes" binding:"required"`
	Views *int   `json:"views" binding:"required"`
	Type  string `json:"type" binding:"required,oneof=Blog Project"`
}

type UpdateStatisticDTO struct {
	Id    int
	Likes int
	Views int
	Type  string
}

type StatisticResponse struct {
	ID        int    `json:"id"`
	Likes     int    `json:"likes"`
	Views     int    `json:"views"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}

type StatisticUpdateResponse struct {
	ID    int    `json:"id"`
	Likes int    `json:"likes"`
	Views int    `json:"views"`
	Type  string `json:"type"`
}

type StatisticDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

func ToStatisticResponse(p Statistic) StatisticResponse {
	return StatisticResponse{
		ID:        p.ID,
		Views:     p.Views,
		Likes:     p.Likes,
		Type:      p.Type,
		CreatedAt: p.CreatedAt.Format("2006-01-02"),
	}
}

func ToStatisticUpdateResponse(p Statistic) StatisticUpdateResponse {
	return StatisticUpdateResponse{
		ID:    p.ID,
		Views: p.Views,
		Likes: p.Likes,
		Type:  p.Type,
	}
}
