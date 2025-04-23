package blog_content_temp_image

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

func (h *handler) GetAll(c *gin.Context) {
	data, err := h.service.GetAllBlogContentTempImgs()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data")
		return
	}
	utils.Success(c, "success get all data", data)
}

func (h *handler) GetBlogContentTempImgById(c *gin.Context) {
	id := c.Param("id")
	data, err := h.service.GetBlogContentTempImgById(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data")
		return
	}
	utils.Success(c, "success get data", data)
}

func (h *handler) ValidateImage(c *gin.Context) (file *multipart.FileHeader, errors []utils.FieldError, err error) {
	validationCheck := []string{"required", "extension", "size"}
	var maxSize int64 = 2 * 1024 * 1024
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}

	// Step 1: Get the file
	image_file, err := c.FormFile("image_file")
	if err != nil && slices.Contains(validationCheck, "required") {
		err_name := fmt.Sprintf("%s is required", "image_file")
		errors = utils.GenerateFieldErrorResponse("image_file", err_name)
		return nil, errors, err
	}

	// Step 2: Validate extension
	errExt := utils.ValidateExtension(image_file.Filename, allowedExtensions)
	if errExt != nil && slices.Contains(validationCheck, "extension") {
		err = fmt.Errorf("validation Error")
		return nil, errExt, err
	}

	// Step 3: Validate size
	if image_file.Size > maxSize && slices.Contains(validationCheck, "size") {
		err = fmt.Errorf("validation Error")
		err_name := fmt.Sprintf("%s exceeds max size", "image_file")
		errors := utils.GenerateFieldErrorResponse("image_file", err_name)
		return nil, errors, err
	}

	return image_file, nil, nil
}

func (h *handler) CreateBlogContentTempImg(c *gin.Context) {
	image_file, errors, err := h.ValidateImage(c)
	if err != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return
	}

	//todo: Upload to minio
	imageRes, err := utils.HandlUploadFile(image_file, "blog")
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to upload file")
		return
	}

	// Validate the struct using validator
	req := CreateBlogContentTempImgRequest{
		ImageUrl:      imageRes.FileURL,
		ImageFileName: imageRes.FileName,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	data, err := h.service.CreateBlogContentTempImg(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to created data")
		return
	}

	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateBlogContentTempImg(c *gin.Context) {
	var req UpdateBlogContentTempImgRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	data, err := h.service.UpdateBlogContentTempImg(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to updated data")
		return
	}

	utils.Success(c, "success updated data", data)
}

func (h *handler) DeleteBlogContentTempImg(c *gin.Context) {
	var req BlogContentTempImgDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	data, err := h.service.DeleteBlogContentTempImg(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to deleted data")
		return
	}
	utils.Success(c, "success deleted data", data)
}
