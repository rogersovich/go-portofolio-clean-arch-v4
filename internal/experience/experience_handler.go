package experience

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
	position := c.DefaultQuery("position", "")
	company_name := c.DefaultQuery("company_name", "")
	work_type := c.DefaultQuery("work_type", "")
	country := c.DefaultQuery("country", "")
	city := c.DefaultQuery("city", "")
	summary_html := c.DefaultQuery("summary_html", "")
	from_date := c.DefaultQuery("from_date", "")
	to_date := c.DefaultQuery("to_date", "")
	is_current := c.DefaultQuery("is_current", "")
	created_at := c.DefaultQuery("created_at", "")

	// Check if the created_at parameter has a value and parse the range
	var createdAtRange []string
	if created_at != "" {
		createdAtRange = strings.Split(created_at, ",")
	}

	var fromDateAtRange []string
	if from_date != "" {
		fromDateAtRange = strings.Split(from_date, ",")
	}

	var toDateAtRange []string
	if to_date != "" {
		toDateAtRange = strings.Split(to_date, ",")
	}

	params := GetAllExperienceParams{
		Page:        page,
		Limit:       limit,
		Sort:        sort,
		Order:       order,
		Position:    position,
		CompanyName: company_name,
		WorkType:    work_type,
		Country:     country,
		City:        city,
		SummaryHTML: summary_html,
		FromDate:    fromDateAtRange,
		ToDate:      toDateAtRange,
		IsCurrent:   is_current,
		CreatedAt:   createdAtRange,
	}

	// Validate the params using the binding tags
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	data, total_records, err := h.service.GetAllExperiences(params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.PaginatedSuccess(c, "success get all data", data, page, limit, total_records)
}

func (h *handler) GetExperienceById(c *gin.Context) {
	// Validate the struct using validator
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid ID")
		return
	}

	data, err := h.service.GetExperienceById(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
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
	comp_image_file, err := c.FormFile("comp_image_file")
	if err != nil && slices.Contains(validationCheck, "required") {
		err_name := fmt.Sprintf("%s is required", "comp_image_file")
		errors = utils.GenerateFieldErrorResponse("comp_image_file", err_name)
		return nil, errors, err
	}

	if slices.Contains(validationCheck, "required") {
		// Step 2: Validate extension
		errExt := utils.ValidateExtension(comp_image_file.Filename, allowedExtensions)
		if errExt != nil && slices.Contains(validationCheck, "extension") {
			err = fmt.Errorf("validation Error")
			return nil, errExt, err
		}

		// Step 3: Validate size
		if comp_image_file.Size > maxSize && slices.Contains(validationCheck, "size") {
			err = fmt.Errorf("validation Error")
			err_name := fmt.Sprintf("%s exceeds max size", "comp_image_file")
			errors := utils.GenerateFieldErrorResponse("comp_image_file", err_name)
			return nil, errors, err
		}
	}

	return comp_image_file, nil, nil
}

func (h *handler) CreateExperience(c *gin.Context) {
	position := c.PostForm("position")
	company_name := c.PostForm("company_name")
	work_type := c.PostForm("work_type")
	country := c.PostForm("country")
	city := c.PostForm("city")
	summary_html := c.PostForm("summary_html")
	from_date := c.PostForm("from_date")
	to_date := c.PostForm("to_date")
	comp_website_url := c.PostForm("comp_website_url")
	is_current := c.PostForm("is_current")

	comp_image_file, errors, err := h.ValidateImage(c, nil)
	if err != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return
	}

	// Validate the struct using validator
	req := CreateExperienceRequest{
		Position:       position,
		CompanyName:    company_name,
		WorkType:       work_type,
		Country:        country,
		City:           &city,
		SummaryHTML:    summary_html,
		FromDate:       from_date,
		ToDate:         &to_date,
		CompImageFile:  comp_image_file,
		CompWebsiteUrl: comp_website_url,
		IsCurrent:      is_current,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	data, err := h.service.CreateExperience(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateExperience(c *gin.Context) {
	// Validate the struct using validator
	id, _ := strconv.Atoi(c.PostForm("id"))

	position := c.PostForm("position")
	company_name := c.PostForm("company_name")
	work_type := c.PostForm("work_type")
	country := c.PostForm("country")
	city := c.PostForm("city")
	summary_html := c.PostForm("summary_html")
	from_date := c.PostForm("from_date")
	to_date := c.PostForm("to_date")
	comp_website_url := c.PostForm("comp_website_url")
	is_current := c.PostForm("is_current")

	validationCheck := []string{"extension", "size"}
	comp_image_file, errors, err := h.ValidateImage(c, validationCheck)
	if err != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return
	}

	req := UpdateExperienceRequest{
		ID:             id,
		Position:       position,
		CompanyName:    company_name,
		WorkType:       work_type,
		Country:        country,
		City:           &city,
		SummaryHTML:    summary_html,
		FromDate:       from_date,
		ToDate:         &to_date,
		CompImageFile:  comp_image_file,
		CompWebsiteUrl: comp_website_url,
		IsCurrent:      is_current,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	err = h.service.UpdateExperience(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success updated data", nil)
}

func (h *handler) DeleteExperience(c *gin.Context) {
	var req ExperienceDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	data, err := h.service.DeleteExperience(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success deleted data", data)
}
