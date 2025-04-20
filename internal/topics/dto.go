package topics

type CreateTopicRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateTopicRequest struct {
	Id   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type TopicResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type TopicUpdateResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TopicDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

func ToTopicResponse(p Topic) TopicResponse {
	return TopicResponse{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt.Format("2006-01-02"),
	}
}

func ToTopicUpdateResponse(p Topic) TopicUpdateResponse {
	return TopicUpdateResponse{
		ID:   p.ID,
		Name: p.Name,
	}
}
