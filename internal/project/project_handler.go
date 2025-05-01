package project

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
	data, err := h.service.GetAllProjects()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data")
		return
	}
	utils.Success(c, "success get all data", data)
}

func (h *handler) GetProjectByIdWithRelations(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid ID")
		return
	}
	data, err := h.service.GetProjectByIdWithRelations(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data")
		return
	}
	utils.Success(c, "success get data", data)
}

func (h *handler) ValidateImage(c *gin.Context, validationCheck []string) (file *multipart.FileHeader, errors []utils.FieldError, err error) {
	if len(validationCheck) == 0 {
		validationCheck = []string{"required", "extension", "size"}
	}
	var maxSize int64 = 2 * 1024 * 1024
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}

	// Step 1: Get the file
	image_file, err := c.FormFile("image_file")
	if err != nil && slices.Contains(validationCheck, "required") {
		err_name := fmt.Sprintf("%s is required", "image_file")
		errors = utils.GenerateFieldErrorResponse("image_file", err_name)
		return nil, errors, err
	}

	if slices.Contains(validationCheck, "required") {
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
	}

	return image_file, nil, nil
}

func (h *handler) CreateProject(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")
	is_published := c.PostForm("is_published") // Y or N
	repository_url := c.PostForm("repository_url")
	summary := c.PostForm("summary")

	var technology_ids []int
	if err := json.Unmarshal([]byte(c.PostForm("technology_ids")), &technology_ids); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid technology_ids format")
		return
	}

	var project_images []string
	if contentImagesParam := c.PostForm("project_images"); contentImagesParam != "" {
		if err := json.Unmarshal([]byte(contentImagesParam), &project_images); err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid project_images format")
			return
		}
	} else {
		project_images = []string{}
	}

	image_file, errors, err := h.ValidateImage(c, nil)
	if err != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return
	}

	// Validate the struct using validator
	req := CreateProjectRequest{
		Title:         title,
		Description:   description,
		ImageFile:     image_file,
		RepositoryUrl: &repository_url,
		Summary:       summary,
		IsPublished:   is_published,
		TechnologyIds: technology_ids,
		ContentImages: project_images,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	data, err := h.service.CreateProject(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateProject(c *gin.Context) {
	// Validate the struct using validator
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid ID")
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	is_published := c.PostForm("is_published") // Y or N
	repository_url := c.PostForm("repository_url")
	summary := c.PostForm("summary")

	var technologyIds []ProjectTechUpdatePayload
	if err := json.Unmarshal([]byte(c.PostForm("technology_ids")), &technologyIds); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid technology_ids format")
		return
	}

	var project_images []string
	if err := json.Unmarshal([]byte(c.PostForm("project_images")), &project_images); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid project_images format")
		return
	}

	validationCheck := []string{"extension", "size"}
	image_file, errors, err := h.ValidateImage(c, validationCheck)
	if err != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return
	}

	// Validate the struct using validator
	req := UpdateProjectRequest{
		Id:            id,
		Title:         title,
		Description:   description,
		ImageFile:     image_file,
		RepositoryUrl: &repository_url,
		Summary:       summary,
		IsPublished:   is_published,
		TechnologyIds: technologyIds,
		ProjectImages: project_images,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	data, err := h.service.UpdateProject(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to updated data")
		return
	}

	utils.Success(c, "success updated data", data)
}

func (h *handler) UpdateProjectStatistic(c *gin.Context) {
	var req ProjectStatisticUpdateRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	data, err := h.service.UpdateProjectStatistic(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to updated data")
		return
	}

	utils.Success(c, "success updated data", data)
}

func (h *handler) DeleteProject(c *gin.Context) {
	var req ProjectDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	data, err := h.service.DeleteProject(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to deleted data")
		return
	}
	utils.Success(c, "success deleted data", data)
}
