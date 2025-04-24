package blog_content_temp_image

type CreateBlogContentTempImgRequest struct {
	ImageUrl      string `json:"image_url" binding:"required"`
	ImageFileName string `json:"image_file_name" binding:"required"`
}

type UpdateBlogContentTempImgRequest struct {
	Id            int    `json:"id" binding:"required"`
	ImageUrl      string `json:"image_url" binding:"required"`
	ImageFileName string `json:"image_file_name" binding:"required"`
}

type BlogContentTempImgResponse struct {
	ID            int    `json:"id"`
	ImageUrl      string `json:"image_url"`
	ImageFileName string `json:"image_file_name"`
	CreatedAt     string `json:"created_at"`
}

type BlogContentTempImgUpdateResponse struct {
	ID            int    `json:"id"`
	ImageUrl      string `json:"image_url"`
	ImageFileName string `json:"image_file_name"`
}

type BlogContentTempImgDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

type CountTempImagesDTO struct {
	ID       int    `json:"id"`
	ImageUrl string `json:"image_url"`
}

type CountTempUpdateImagesDTO struct {
	ID          int    `json:"id"`
	ImageUrl    string `json:"image_url"`
	ImageOldUrl string `json:"image_old_url"`
}

func ToBlogContentTempImgResponse(p BlogContentTempImages) BlogContentTempImgResponse {
	return BlogContentTempImgResponse{
		ID:            p.ID,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		CreatedAt:     p.CreatedAt.Format("2006-01-02"),
	}
}

func ToBlogContentTempImgUpdateResponse(p BlogContentTempImages) BlogContentTempImgUpdateResponse {
	return BlogContentTempImgUpdateResponse{
		ID:            p.ID,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
	}
}
