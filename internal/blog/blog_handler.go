package blog

import (
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
		utils.Error(c, http.StatusInternalServerError, "failed to get data", err)
		return
	}
	utils.Success(c, "success get all data", data)
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
	topic_ids := c.PostFormArray("topic_ids[]")

	author_id_int, err := strconv.Atoi(author_id)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid author_id", err)
		return
	}

	topic_ids_validated, err := utils.ValidateFormArrayString(topic_ids, "topic_ids", true)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid topic_ids", err)
		return
	}

	image_file, errors, err := h.ValidateBanner(c, nil)
	if err != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return
	}

	// Validate the struct using validator
	req := CreateBlogRequest{
		TopicIds:        topic_ids_validated,
		AuthorID:        author_id_int,
		Title:           title,
		DescriptionHTML: description,
		BannerFile:      image_file,
		Summary:         summary,
		IsPublished:     is_published,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	utils.PrintJSON(req)

	// data, err := h.service.CreateBlog(req)
	// if err != nil {
	// 	utils.Error(c, http.StatusInternalServerError, "failed to created data", err)
	// 	return
	// }

	utils.Success(c, "success get data", nil)
}
