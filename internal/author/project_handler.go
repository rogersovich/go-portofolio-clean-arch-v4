package author

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
	data, err := h.service.GetAllAuthors()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data", err)
		return
	}
	utils.Success(c, "success get all data", data)
}

func (h *handler) GetAuthorById(c *gin.Context) {
	id := c.Param("id")
	data, err := h.service.GetAuthorById(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data", err)
		return
	}
	utils.Success(c, "success get data", data)
}

func (h *handler) ValidateAvatar(c *gin.Context) (file *multipart.FileHeader, errors []utils.FieldError, err error) {
	validationCheck := []string{"required", "extension", "size"}
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

	return avatar_file, nil, nil
}

func (h *handler) CreateAuthor(c *gin.Context) {
	name := c.PostForm("name")

	avatar_file, errors, err := h.ValidateAvatar(c)
	if err != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return
	}

	avatarRes, err := utils.HandlUploadFile(avatar_file, "author")
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to upload file", err)
		return
	}

	// Validate the struct using validator
	req := CreateAuthorRequest{
		Name:           name,
		AvatarUrl:      avatarRes.FileURL,
		AvatarFileName: avatarRes.FileName,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	data, err := h.service.CreateAuthor(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to created data", err)
		return
	}

	utils.PrintJSON(data)

	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateAuthor(c *gin.Context) {
	// Validate the struct using validator
	id, _ := strconv.Atoi(c.PostForm("id"))

	name := c.PostForm("name")
	req := UpdateAuthorRequest{
		Id:   id,
		Name: name,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	author, err := h.service.GetAuthorById(strconv.Itoa(id))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Data not found", err)
		return
	}

	// set oldPath
	oldPath := ""
	if author.AvatarFileName != "" {
		oldPath = author.AvatarFileName
	}

	// 2. Get new file (if uploaded)
	_, err = c.FormFile("avatar_file")
	var newFileURL string
	var newFileName string

	if err == nil {
		avatar_file, errors, err := h.ValidateAvatar(c)
		if err != nil {
			utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
			return
		}

		avatarRes, err := utils.HandlUploadFile(avatar_file, "author")
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to upload file", err)
			return
		}

		newFileURL = avatarRes.FileURL
		newFileName = avatarRes.FileName
	} else {
		newFileURL = author.AvatarUrl // keep existing if not updated
		newFileName = author.AvatarFileName
	}

	// Validate the struct using validator
	payload := UpdateAuthorDTO{
		Id:             uint(id),
		Name:           name,
		AvatarUrl:      newFileURL,
		AvatarFileName: newFileName,
	}

	data, err := h.service.UpdateAuthor(payload, oldPath, newFileName)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to updated data", err)
		return
	}

	utils.Success(c, "success updated data", data)
}

func (h *handler) DeleteAuthor(c *gin.Context) {
	var req AuthorDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	data, err := h.service.DeleteAuthor(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to deleted data", err)
		return
	}
	utils.Success(c, "success deleted data", data)
}
