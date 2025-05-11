package reading_time

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
	min_minutes := c.DefaultQuery("min_minutes", "")
	max_minutes := c.DefaultQuery("max_minutes", "")
	min_estimates := c.DefaultQuery("min_estimates", "")
	max_estimates := c.DefaultQuery("max_estimates", "")
	created_at := c.DefaultQuery("created_at", "")

	// Check if the created_at parameter has a value and parse the range
	var createdAtRange []string
	if created_at != "" {
		createdAtRange = strings.Split(created_at, ",")
	}

	params := GetAllReadingTimeParams{
		Page:         page,
		Limit:        limit,
		Sort:         sort,
		Order:        order,
		MinMinutes:   min_minutes,
		MaxMinutes:   max_minutes,
		MinEstimates: min_estimates,
		MaxEstimates: max_estimates,
		CreatedAt:    createdAtRange,
	}

	// Validate the params using the binding tags
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	data, total_records, err := h.service.GetAllReadingTimes(params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.PaginatedSuccess(c, "success get all data", data, page, limit, total_records)
}

func (h *handler) GetReadingTimeById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid ID")
		return
	}
	data, err := h.service.GetReadingTimeById(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get data", data)
}

func (h *handler) CreateReadingTime(c *gin.Context) {
	// Validate the struct using validator
	var req CreateReadingTimeRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	data, err := h.service.CreateReadingTime(req, nil)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateReadingTime(c *gin.Context) {
	var req UpdateReadingTimeRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	err := h.service.UpdateReadingTime(req, nil)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success updated data", nil)
}

func (h *handler) DeleteReadingTime(c *gin.Context) {
	var req ReadingTimeDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	data, err := h.service.DeleteReadingTime(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success deleted data", data)
}
