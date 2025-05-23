package topic

type CreateTopicRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateTopicRequest struct {
	ID   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type TopicResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type TopicHasCheckResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TopicDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

type TopicCheckIdsRequest struct {
	Ids []int `json:"ids" binding:"required,dive,gt=0"`
}

type GetAllTopicParams struct {
	Limit     int `binding:"required"`
	Page      int `binding:"required"`
	Sort      string
	Order     string
	Name      string
	CreatedAt []string
}

func ToTopicResponse(p Topic) TopicResponse {
	return TopicResponse{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt.Format("2006-01-02"),
	}
}
