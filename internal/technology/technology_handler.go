package technology

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

func (h *handler) GetAll(c *gin.Context) {
	page := utils.GetQueryParamInt(c, "page", 1) // Default to page 1
	limit := utils.GetQueryParamInt(c, "limit", 10)
	//? Sort and order
	sort := c.DefaultQuery("sort", "ASC")
	order := c.DefaultQuery("order", "id")
	//? Filters
	name := c.DefaultQuery("name", "")
	description_html := c.DefaultQuery("description_html", "")
	is_major := c.DefaultQuery("is_major", "")
	created_at := c.DefaultQuery("created_at", "")

	// Check if the created_at parameter has a value and parse the range
	var createdAtRange []string
	if created_at != "" {
		createdAtRange = strings.Split(created_at, ",")
	}

	params := GetAllTechnologyParams{
		Page:            page,
		Limit:           limit,
		Sort:            sort,
		Order:           order,
		Name:            name,
		DescriptionHTML: description_html,
		IsMajor:         is_major,
		CreatedAt:       createdAtRange,
	}

	// Validate the params using the binding tags
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	data, total_records, err := h.service.GetAllTechnologies(params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.PaginatedSuccess(c, "success get all data", data, page, limit, total_records)
}

func (h *handler) GetTechnologyById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid ID")
		return
	}
	data, err := h.service.GetTechnologyById(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get data", data)
}

func (h *handler) ValidateLogo(c *gin.Context, validationCheck []string) (file *multipart.FileHeader, errors []utils.FieldError, err error) {
	if len(validationCheck) == 0 {
		validationCheck = []string{"required", "extension", "size"}
	}
	var maxSize int64 = 2 * 1024 * 1024
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}

	// Step 1: Get the file
	logo_file, err := c.FormFile("logo_file")
	if err != nil && slices.Contains(validationCheck, "required") {
		err_name := fmt.Sprintf("%s is required", "logo_file")
		errors = utils.GenerateFieldErrorResponse("logo_file", err_name)
		return nil, errors, err
	}

	if slices.Contains(validationCheck, "required") {
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
	}

	return logo_file, nil, nil
}

func (h *handler) CreateTechnology(c *gin.Context) {
	name := c.PostForm("name")
	description_html := c.PostForm("description_html")
	is_major := c.PostForm("is_major")
	link := c.PostForm("link")

	logo_file, errors, err := h.ValidateLogo(c, nil)
	if err != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return
	}

	// Validate the struct using validator
	req := CreateTechnologyRequest{
		Name:            name,
		DescriptionHTML: description_html,
		LogoFile:        logo_file,
		IsMajor:         is_major,
		Link:            &link,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	data, err := h.service.CreateTechnology(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateTechnology(c *gin.Context) {
	// Validate the struct using validator
	id, _ := strconv.Atoi(c.PostForm("id"))

	name := c.PostForm("name")
	description_html := c.PostForm("description_html")
	is_major := c.PostForm("is_major")
	link := c.PostForm("link")

	validationCheck := []string{"extension", "size"}
	logo_file, errors, err := h.ValidateLogo(c, validationCheck)
	if err != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return
	}

	req := UpdateTechnologyRequest{
		ID:              id,
		Name:            name,
		DescriptionHTML: description_html,
		LogoFile:        logo_file,
		IsMajor:         is_major,
		Link:            &link,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	err = h.service.UpdateTechnology(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success updated data", nil)
}

func (h *handler) DeleteTechnology(c *gin.Context) {
	var req TechnologyDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	data, err := h.service.DeleteTechnology(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success deleted data", data)
}
