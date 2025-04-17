package technology

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
	data, err := h.service.GetAllTechnologies()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data", err)
		return
	}
	utils.Success(c, "success get all data", data)
}

func (h *handler) GetTechnologyById(c *gin.Context) {
	id := c.Param("id")
	data, err := h.service.GetTechnologyById(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data", err)
		return
	}
	utils.Success(c, "success get data", data)
}

func (h *handler) ValidateLogo(c *gin.Context) (file *multipart.FileHeader, errors []utils.FieldError, err error) {
	validationCheck := []string{"required", "extension", "size"}
	var maxSize int64 = 2 * 1024 * 1024
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}

	// Step 1: Get the file
	logo_file, err := c.FormFile("logo_file")
	if err != nil && slices.Contains(validationCheck, "required") {
		err_name := fmt.Sprintf("%s is required", "logo_file")
		errors = utils.GenerateFieldErrorResponse("logo_file", err_name)
		return nil, errors, err
	}

	// Step 2: Validate extension
	errExt := utils.ValidateExtension(logo_file.Filename, allowedExtensions)
	if errExt != nil && slices.Contains(validationCheck, "extension") {
		err = fmt.Errorf("validation Error")
		return nil, errExt, err
	}

	// Step 3: Validate size
	if logo_file.Size > maxSize && slices.Contains(validationCheck, "size") {
		err = fmt.Errorf("validation Error")
		err_name := fmt.Sprintf("%s exceeds max size", "logo_file")
		errors := utils.GenerateFieldErrorResponse("logo_file", err_name)
		return nil, errors, err
	}

	return logo_file, nil, nil
}

func (h *handler) CreateTechnology(c *gin.Context) {
	name := c.PostForm("name")
	description_html := c.PostForm("description_html")
	is_major := c.PostForm("is_major")

	logo_file, errors, err := h.ValidateLogo(c)
	if err != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return
	}

	logoRes, err := utils.HandlUploadFile(logo_file, "technology")
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to upload file", err)
		return
	}

	// Validate the struct using validator
	req := CreateTechnologyRequest{
		Name:            name,
		DescriptionHTML: description_html,
		LogoUrl:         logoRes.FileURL,
		LogoFileName:    logoRes.FileName,
		IsMajor:         is_major,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	data, err := h.service.CreateTechnology(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to created data", err)
		return
	}

	utils.PrintJSON(data)

	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateTechnology(c *gin.Context) {
	// Validate the struct using validator
	id, _ := strconv.Atoi(c.PostForm("id"))

	name := c.PostForm("name")
	description_html := c.PostForm("description_html")
	is_major := c.PostForm("is_major")
	req := UpdateTechnologyRequest{
		Id:              id,
		Name:            name,
		DescriptionHTML: description_html,
		IsMajor:         is_major,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	technology, err := h.service.GetTechnologyById(strconv.Itoa(id))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Data not found", err)
		return
	}

	// set oldPath
	oldPath := ""
	if technology.LogoFileName != "" {
		oldPath = technology.LogoFileName
	}

	// 2. Get new file (if uploaded)
	_, err = c.FormFile("logo_file")
	var newFileURL string
	var newFileName string

	if err == nil {
		logo_file, errors, err := h.ValidateLogo(c)
		if err != nil {
			utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
			return
		}

		logoRes, err := utils.HandlUploadFile(logo_file, "technology")
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to upload file", err)
			return
		}

		newFileURL = logoRes.FileURL
		newFileName = logoRes.FileName
	} else {
		newFileURL = technology.LogoUrl // keep existing if not updated
		newFileName = technology.LogoFileName
	}

	// Validate the struct using validator
	payload := UpdateTechnologyDTO{
		Id:              uint(id),
		Name:            name,
		DescriptionHTML: description_html,
		LogoUrl:         newFileURL,
		LogoFileName:    newFileName,
		IsMajor:         is_major,
	}

	data, err := h.service.UpdateTechnology(payload, oldPath, newFileName)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to updated data", err)
		return
	}

	utils.Success(c, "success updated data", data)
}

func (h *handler) DeleteTechnology(c *gin.Context) {
	var req TechnologyDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	data, err := h.service.DeleteTechnology(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to deleted data", err)
		return
	}
	utils.Success(c, "success deleted data", data)
}
