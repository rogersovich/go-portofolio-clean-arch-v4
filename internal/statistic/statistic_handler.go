package statistic

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

func (h *handler) GetAll(c *gin.Context) {
	data, err := h.service.GetAllStatistics()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data")
		return
	}
	utils.Success(c, "success get all data", data)
}

func (h *handler) GetStatisticById(c *gin.Context) {
	id := c.Param("id")
	data, err := h.service.GetStatisticById(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data")
		return
	}
	utils.Success(c, "success get data", data)
}

func (h *handler) CreateStatistic(c *gin.Context) {
	// Validate the struct using validator
	var req CreateStatisticRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	data, err := h.service.CreateStatistic(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to created data")
		return
	}

	utils.PrintJSON(data)

	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateStatistic(c *gin.Context) {
	var req UpdateStatisticRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	data, err := h.service.UpdateStatistic(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to updated data")
		return
	}

	utils.Success(c, "success updated data", data)
}

func (h *handler) DeleteStatistic(c *gin.Context) {
	var req StatisticDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	data, err := h.service.DeleteStatistic(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to deleted data")
		return
	}
	utils.Success(c, "success deleted data", data)
}
