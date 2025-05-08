package author

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
	created_at := c.DefaultQuery("created_at", "")

	// Check if the created_at parameter has a value and parse the range
	var createdAtRange []string
	if created_at != "" {
		created_at = strings.Trim(created_at, "[]")
		createdAtRange = strings.Split(created_at, ",")
	}

	params := GetAllAuthorParams{
		Page:      page,
		Limit:     limit,
		Sort:      sort,
		Order:     order,
		Name:      name,
		CreatedAt: createdAtRange,
	}

	// Validate the params using the binding tags
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	data, total_records, err := h.service.GetAllAuthors(params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.PaginatedSuccess(c, "success get all data", data, page, limit, int(total_records))
}

func (h *handler) GetAuthorById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid ID")
		return
	}
	data, err := h.service.GetAuthorById(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get data", data)
}

func (h *handler) ValidateAvatar(c *gin.Context, validationCheck []string) (file *multipart.FileHeader, errors []utils.FieldError, err error) {
	if len(validationCheck) == 0 {
		validationCheck = []string{"required", "extension", "size"}
	}
	var maxSize int64 = 2 * 1024 * 1024
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}

	// Step 1: Get the file
	avatar_file, err := c.FormFile("avatar_file")
	if err != nil && slices.Contains(validationCheck, "required") {
		err_name := fmt.Sprintf("%s is required", "avatar_file")
		errors = utils.GenerateFieldErrorResponse("avatar_file", err_name)
		// utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return nil, errors, err
	}

	if slices.Contains(validationCheck, "required") {
		// Step 2: Validate extension
		errExt := utils.ValidateExtension(avatar_file.Filename, allowedExtensions)
		if errExt != nil && slices.Contains(validationCheck, "extension") {
			err = fmt.Errorf("validation Error")
			// utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", errExt)
			return nil, errExt, err
		}

		// Step 3: Validate size
		if avatar_file.Size > maxSize && slices.Contains(validationCheck, "size") {
			err = fmt.Errorf("validation Error")
			err_name := fmt.Sprintf("%s exceeds max size", "avatar_file")
			errors := utils.GenerateFieldErrorResponse("avatar_file", err_name)
			// utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", errors)
			return nil, errors, err
		}
	}

	return avatar_file, nil, nil
}

func (h *handler) CreateAuthor(c *gin.Context) {
	name := c.PostForm("name")

	avatar_file, errors, err := h.ValidateAvatar(c, nil)
	if err != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return
	}

	// Validate the struct using validator
	req := CreateAuthorRequest{
		Name:       name,
		AvatarFile: avatar_file,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	data, err := h.service.CreateAuthor(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateAuthor(c *gin.Context) {
	// Validate the struct using validator
	id, _ := strconv.Atoi(c.PostForm("id"))

	name := c.PostForm("name")

	validationCheck := []string{"extension", "size"}
	avatar_file, errors, err := h.ValidateAvatar(c, validationCheck)
	if err != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return
	}

	req := UpdateAuthorRequest{
		ID:         id,
		Name:       name,
		AvatarFile: avatar_file,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	err = h.service.UpdateAuthor(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success updated data", nil)
}

func (h *handler) DeleteAuthor(c *gin.Context) {
	var req AuthorDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	err := h.service.DeleteAuthor(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success deleted data", nil)
}
