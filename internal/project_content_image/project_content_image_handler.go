package project_content_image

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
	datas, err := h.service.GetAllProjectContentImages()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data")
		return
	}
	utils.Success(c, "success get all data", datas)
}

func (h *handler) GetProjectContentImageById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id type is wrong")
		return
	}
	data, err := h.service.GetProjectContentImageById(id)
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

func (h *handler) CreateProjectContentImage(c *gin.Context) {
	image_file, errors, err := h.ValidateImage(c)
	if err != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return
	}

	imageRes, err := utils.HandlUploadFile(image_file, "project")
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to upload file")
		return
	}

	// Validate the struct using validator
	req := CreateProjectContentImageRequest{
		ImageUrl:      imageRes.FileURL,
		ImageFileName: imageRes.FileName,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	data, err := h.service.CreateProjectContentImage(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to created data")
		return
	}

	utils.PrintJSON(data)

	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateProjectContentImage(c *gin.Context) {
	// Validate the struct using validator
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id type is wrong")
		return
	}
	project_id, err := strconv.Atoi(c.PostForm("project_id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "project_id type is wrong")
		return
	}

	req := UpdateProjectContentImageRequest{
		Id:        id,
		ProjectID: &project_id,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	data, err := h.service.GetProjectContentImageById(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Data not found")
		return
	}

	// set oldPath
	oldPath := ""
	if data.ImageFileName != "" {
		oldPath = data.ImageFileName
	}

	// 2. Get new file (if uploaded)
	_, err = c.FormFile("image_file")
	var newFileURL string
	var newFileName string

	if err == nil {
		image_file, errors, err := h.ValidateImage(c)
		if err != nil {
			utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
			return
		}

		imageRes, err := utils.HandlUploadFile(image_file, "project")
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to upload file")
			return
		}

		newFileURL = imageRes.FileURL
		newFileName = imageRes.FileName
	} else {
		newFileURL = data.ImageUrl // keep existing if not updated
		newFileName = data.ImageFileName
	}

	// Validate the struct using validator
	payload := UpdateProjectContentImageDTO{
		Id:            id,
		ProjectID:     &project_id,
		ImageUrl:      newFileURL,
		ImageFileName: newFileName,
	}

	result, err := h.service.UpdateProjectContentImage(payload, oldPath, newFileName)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to updated data")
		return
	}

	utils.Success(c, "success updated data", result)
}

func (h *handler) DeleteProjectContentImage(c *gin.Context) {
	var req ProjectContentImageDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	data, err := h.service.DeleteProjectContentImage(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to deleted data")
		return
	}
	utils.Success(c, "success deleted data", data)
}
