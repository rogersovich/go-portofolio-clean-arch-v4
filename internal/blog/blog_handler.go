package blog

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

func (h *handler) GetAll(c *gin.Context) {
	data, err := h.service.GetAllBlogs()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data")
		return
	}
	utils.Success(c, "success get all data", data)
}

func (h *handler) GetBlogByIdWithRelations(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid ID")
		return
	}
	data, err := h.service.GetBlogByIdWithRelations(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get data", data)
}

func (h *handler) ValidateBanner(c *gin.Context, validationCheck []string) (file *multipart.FileHeader, errors []utils.FieldError, err error) {
	if len(validationCheck) == 0 {
		validationCheck = []string{"required", "extension", "size"}
	}
	var maxSize int64 = 2 * 1024 * 1024
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}

	// Step 1: Get the file
	banner_file, err := c.FormFile("banner_file")
	if err != nil && slices.Contains(validationCheck, "required") {
		err_name := fmt.Sprintf("%s is required", "banner_file")
		errors = utils.GenerateFieldErrorResponse("banner_file", err_name)
		return nil, errors, err
	}

	if slices.Contains(validationCheck, "required") {
		// Step 2: Validate extension
		errExt := utils.ValidateExtension(banner_file.Filename, allowedExtensions)
		if errExt != nil && slices.Contains(validationCheck, "extension") {
			err = fmt.Errorf("validation Error")
			return nil, errExt, err
		}

		// Step 3: Validate size
		if banner_file.Size > maxSize && slices.Contains(validationCheck, "size") {
			err = fmt.Errorf("validation Error")
			err_name := fmt.Sprintf("%s exceeds max size", "banner_file")
			errors := utils.GenerateFieldErrorResponse("banner_file", err_name)
			return nil, errors, err
		}
	}

	return banner_file, nil, nil
}

func (h *handler) CreateBlog(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")
	is_published := c.PostForm("is_published") // Y or N
	summary := c.PostForm("summary")
	author_id := c.PostForm("author_id")
	slug := c.PostForm("slug")

	author_id_int, err := strconv.Atoi(author_id)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid author_id")
		return
	}

	var topic_ids []int
	if err := json.Unmarshal([]byte(c.PostForm("topic_ids")), &topic_ids); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid topic_ids format")
		return
	}

	var content_images []string

	// Check if "content_images" field is provided in the form
	if contentImagesParam := c.PostForm("content_images"); contentImagesParam != "" {
		// Unmarshal if the content_images field is provided
		if err := json.Unmarshal([]byte(contentImagesParam), &content_images); err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid content_images format")
			return
		}
	} else {
		// If content_images is not provided, assign an empty slice or leave it nil
		content_images = []string{}
	}

	image_file, errors, err := h.ValidateBanner(c, nil)
	if err != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return
	}

	// Validate the struct using validator
	req := CreateBlogRequest{
		TopicIds:        topic_ids,
		AuthorID:        author_id_int,
		Title:           title,
		DescriptionHTML: description,
		BannerFile:      image_file,
		Summary:         summary,
		IsPublished:     is_published,
		ContentImages:   content_images,
		Slug:            slug,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	data, err := h.service.CreateBlog(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateBlog(c *gin.Context) {
	// Validate the struct using validator
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid ID")
		return
	}
	// Validate the struct using validator
	author_id, err := strconv.Atoi(c.PostForm("author_id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid Author ID")
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	is_published := c.PostForm("is_published") // Y or N
	summary := c.PostForm("summary")
	slug := c.PostForm("slug")

	var topic_ids []UpdateBlogTopicDTO
	if err := json.Unmarshal([]byte(c.PostForm("topic_ids")), &topic_ids); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid topic_ids format")
		return
	}

	var content_images []string
	if err := json.Unmarshal([]byte(c.PostForm("content_images")), &content_images); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid content_images format")
		return
	}

	validationCheck := []string{"extension", "size"}
	banner_file, errors, err := h.ValidateBanner(c, validationCheck)
	if err != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return
	}

	// Validate the struct using validator
	req := UpdateBlogRequest{
		ID:              id,
		Title:           title,
		DescriptionHTML: description,
		BannerFile:      banner_file,
		Summary:         summary,
		IsPublished:     is_published,
		TopicIds:        topic_ids,
		ContentImages:   content_images,
		AuthorID:        author_id,
		Slug:            slug,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	data, err := h.service.UpdateBlog(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success get data", data)
}

func (h *handler) DeleteBlog(c *gin.Context) {
	var req BlogDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	data, err := h.service.DeleteBlog(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to deleted data")
		return
	}
	utils.Success(c, "success deleted data", data)
}

func (h *handler) ChangeStatusBlog(c *gin.Context) {
	var req BlogChangeStatusRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	data, err := h.service.ChangeStatusBlog(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success change status", data)
}
