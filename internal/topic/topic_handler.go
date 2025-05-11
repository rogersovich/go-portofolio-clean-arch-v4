package topic

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
	name := c.DefaultQuery("name", "")
	created_at := c.DefaultQuery("created_at", "")

	// Check if the created_at parameter has a value and parse the range
	var createdAtRange []string
	if created_at != "" {
		createdAtRange = strings.Split(created_at, ",")
	}

	params := GetAllTopicParams{
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

	data, total_records, err := h.service.GetAllTopics(params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.PaginatedSuccess(c, "success get all data", data, page, limit, total_records)
}

func (h *handler) GetTopicById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid ID")
		return
	}
	data, err := h.service.GetTopicById(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get data", data)
}

func (h *handler) CreateTopic(c *gin.Context) {
	// Validate the struct using validator
	var req CreateTopicRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	data, err := h.service.CreateTopic(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateTopic(c *gin.Context) {
	var req UpdateTopicRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	err := h.service.UpdateTopic(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success updated data", nil)
}

func (h *handler) DeleteTopic(c *gin.Context) {
	var req TopicDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	data, err := h.service.DeleteTopic(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success deleted data", data)
}

func (h *handler) CheckTopicIds(c *gin.Context) {
	var req TopicCheckIdsRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	data, err := h.service.CheckTopicIds(req.Ids)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success fetch data", data)
}
