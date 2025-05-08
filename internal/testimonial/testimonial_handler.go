package testimonial

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
	role := c.DefaultQuery("role", "")
	working_at := c.DefaultQuery("working_at", "")
	is_used := c.DefaultQuery("is_used", "")
	created_at := c.DefaultQuery("created_at", "")

	// Check if the created_at parameter has a value and parse the range
	var createdAtRange []string
	if created_at != "" {
		created_at = strings.Trim(created_at, "[]")
		createdAtRange = strings.Split(created_at, ",")
	}

	params := GetAllTestimonialParams{
		Page:      page,
		Limit:     limit,
		Sort:      sort,
		Order:     order,
		Name:      name,
		Role:      role,
		WorkingAt: working_at,
		IsUsed:    is_used,
		CreatedAt: createdAtRange,
	}

	// Validate the params using the binding tags
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	data, total_records, err := h.service.GetAllTestimonials(params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.PaginatedSuccess(c, "success get all data", data, page, limit, total_records)
}

func (h *handler) GetTestimonialById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid ID")
		return
	}
	data, err := h.service.GetTestimonialById(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get data", data)
}

func (h *handler) CreateTestimonial(c *gin.Context) {
	// Validate the struct using validator
	var req CreateTestimonialRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	data, err := h.service.CreateTestimonial(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateTestimonial(c *gin.Context) {
	var req UpdateTestimonialRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	err := h.service.UpdateTestimonial(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success updated data", nil)
}

func (h *handler) DeleteTestimonial(c *gin.Context) {
	var req TestimonialDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	data, err := h.service.DeleteTestimonial(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success deleted data", data)
}

func (h *handler) ChangeStatusTestimonial(c *gin.Context) {
	var req TestimonialChangeStatusRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	err := h.service.ChangeStatusTestimonial(req.ID, req.IsUsed)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success change status", nil)

}
func (h *handler) ChangeMultiStatusTestimonial(c *gin.Context) {
	var req TestimonialChangeMultiStatusRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	err := h.service.ChangeMultiStatusTestimonial(req.IDs, req.IsUsed)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success change status", nil)

}
