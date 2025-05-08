package user

import (
	"net/http"
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
	username := c.DefaultQuery("username", "")
	email := c.DefaultQuery("email", "")
	created_at := c.DefaultQuery("created_at", "")

	// Check if the created_at parameter has a value and parse the range
	var createdAtRange []string
	if created_at != "" {
		created_at = strings.Trim(created_at, "[]")
		createdAtRange = strings.Split(created_at, ",")
	}

	params := GetAllUserParams{
		Page:      page,
		Limit:     limit,
		Sort:      sort,
		Order:     order,
		Username:  username,
		Email:     email,
		CreatedAt: createdAtRange,
	}

	// Validate the params using the binding tags
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	data, total_records, err := h.service.GetAllUsers(params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.PaginatedSuccess(c, "success get all data", data, page, limit, total_records)
}

func (h *handler) GetUserById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid ID")
		return
	}
	data, err := h.service.GetUserById(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateUser(c *gin.Context) {
	// Validate the struct using validator
	var req UserUpdateRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	payload := User{
		ID:       req.ID,
		Username: req.Username,
		Email:    req.Email,
	}

	data, err := h.service.UpdateUser(payload)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success updated data", data)
}

func (h *handler) DeleteUser(c *gin.Context) {
	// Validate the struct using validator
	var req UserDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	err := h.service.DeleteUser(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success deleted data", nil)
}
