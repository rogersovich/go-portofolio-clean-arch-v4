package about

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
	data, err := h.service.GetAllAbouts()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data")
		return
	}
	utils.Success(c, "success get all data", data)
}

func (h *handler) GetAboutById(c *gin.Context) {
	id := c.Param("id")
	data, err := h.service.GetAboutById(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data")
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
		return nil, errors, err
	}

	// Step 2: Validate extension
	errExt := utils.ValidateExtension(avatar_file.Filename, allowedExtensions)
	if errExt != nil && slices.Contains(validationCheck, "extension") {
		err = fmt.Errorf("validation Error")
		return nil, errExt, err
	}

	// Step 3: Validate size
	if avatar_file.Size > maxSize && slices.Contains(validationCheck, "size") {
		err = fmt.Errorf("validation Error")
		err_name := fmt.Sprintf("%s exceeds max size", "avatar_file")
		errors := utils.GenerateFieldErrorResponse("avatar_file", err_name)
		return nil, errors, err
	}

	return avatar_file, nil, nil
}

func (h *handler) CreateAbout(c *gin.Context) {
	title := c.PostForm("title")
	description_html := c.PostForm("description_html")

	avatar_file, errors, err := h.ValidateAvatar(c)
	if err != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, err.Error(), errors)
		return
	}

	avatarRes, err := utils.HandlUploadFile(avatar_file, "about")
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to upload file")
		return
	}

	// Validate the struct using validator
	req := CreateAboutRequest{
		Title:           title,
		DescriptionHTML: description_html,
		AvatarUrl:       avatarRes.FileURL,
		AvatarFileName:  avatarRes.FileName,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	data, err := h.service.CreateAbout(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to created data")
		return
	}

	utils.PrintJSON(data)

	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateAbout(c *gin.Context) {
	// Validate the struct using validator
	id, _ := strconv.Atoi(c.PostForm("id"))

	title := c.PostForm("title")
	description_html := c.PostForm("description_html")
	req := UpdateAboutRequest{
		Id:              id,
		Title:           title,
		DescriptionHTML: description_html,
	}

	if verr := utils.ValidateRequest(&req); verr != nil {
		utils.ErrorValidation(c, http.StatusBadRequest, "Validation Error", verr)
		return
	}

	about, err := h.service.GetAboutById(strconv.Itoa(id))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Data not found")
		return
	}

	// set oldPath
	oldPath := ""
	if about.AvatarFileName != "" {
		oldPath = about.AvatarFileName
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

		avatarRes, err := utils.HandlUploadFile(avatar_file, "about")
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to upload file")
			return
		}

		newFileURL = avatarRes.FileURL
		newFileName = avatarRes.FileName
	} else {
		newFileURL = about.AvatarUrl // keep existing if not updated
		newFileName = about.AvatarFileName
	}

	// Validate the struct using validator
	payload := UpdateAboutDTO{
		Id:              uint(id),
		Title:           title,
		DescriptionHTML: description_html,
		AvatarUrl:       newFileURL,
		AvatarFileName:  newFileName,
	}

	data, err := h.service.UpdateAbout(payload, oldPath, newFileName)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to updated data")
		return
	}

	utils.Success(c, "success updated data", data)
}

func (h *handler) DeleteAbout(c *gin.Context) {
	var req AboutDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	data, err := h.service.DeleteAbout(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to deleted data")
		return
	}
	utils.Success(c, "success deleted data", data)
}
