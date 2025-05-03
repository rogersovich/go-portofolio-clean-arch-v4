package testimonial

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

func (h *handler) GetAll(c *gin.Context) {
	data, err := h.service.GetAllTestimonials()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get all data", data)
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
