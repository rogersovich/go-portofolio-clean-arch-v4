package blog_topic

type CreateBlogTopicRequest struct {
	BlogID  int `json:"blog_id" binding:"required"`
	TopicID int `json:"topic_id" binding:"required"`
}

type UpdateBlogTopicRequest struct {
	Id      int `json:"id" binding:"required"`
	BlogID  int `json:"blog_id" binding:"required"`
	TopicID int `json:"topic_id" binding:"required"`
}

type BlogTopicResponse struct {
	ID        int    `json:"id"`
	BlogID    int    `json:"blog_id"`
	TopicID   int    `json:"topic_id"`
	CreatedAt string `json:"created_at"`
}

type BlogTopicUpdateResponse struct {
	ID      int `json:"id"`
	BlogID  int `json:"blog_id"`
	TopicID int `json:"topic_id"`
}

type BlogTopicDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

func ToBlogTopicResponse(p BlogTopic) BlogTopicResponse {
	return BlogTopicResponse{
		ID:        p.ID,
		BlogID:    p.BlogID,
		TopicID:   p.TopicID,
		CreatedAt: p.CreatedAt.Format("2006-01-02"),
	}
}

func ToBlogTopicUpdateResponse(p BlogTopic) BlogTopicUpdateResponse {
	return BlogTopicUpdateResponse{
		ID:      p.ID,
		BlogID:  p.BlogID,
		TopicID: p.TopicID,
	}
}
