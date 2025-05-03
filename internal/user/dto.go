package user

type UserUpdateRequest struct {
	ID       int    `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type UserDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

func ToUserResponse(p User) UserResponse {
	return UserResponse{
		ID:        p.ID,
		Username:  p.Username,
		Email:     p.Email,
		CreatedAt: p.CreatedAt.Format("2006-01-02"),
	}
}
