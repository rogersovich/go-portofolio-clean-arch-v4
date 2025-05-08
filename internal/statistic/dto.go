package statistic

type CreateStatisticRequest struct {
	Likes *int   `json:"likes" binding:"required"`
	Views *int   `json:"views" binding:"required"`
	Type  string `json:"type" binding:"required,oneof=Blog Project"`
}

type UpdateStatisticRequest struct {
	ID    int    `json:"id" binding:"required"`
	Likes *int   `json:"likes" binding:"required"`
	Views *int   `json:"views" binding:"required"`
	Type  string `json:"type" binding:"required,oneof=Blog Project"`
}

type UpdateStatisticDTO struct {
	ID    int
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

type StatisticDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

type GetAllStatisticParams struct {
	Limit     int `binding:"required"`
	Page      int `binding:"required"`
	Sort      string
	Order     string
	Type      string
	MinLikes  string
	MaxLikes  string
	MinViews  string
	MaxViews  string
	CreatedAt []string
}

func ToStatisticResponse(p Statistic) StatisticResponse {
	return StatisticResponse{
		ID:        p.ID,
		Views:     *p.Views,
		Likes:     *p.Likes,
		Type:      p.Type,
		CreatedAt: p.CreatedAt.Format("2006-01-02"),
	}
}
