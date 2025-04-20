package topics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

func (h *handler) GetAll(c *gin.Context) {
	data, err := h.service.GetAllTopics()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data")
		return
	}
	utils.Success(c, "success get all data", data)
}

func (h *handler) GetTopicById(c *gin.Context) {
	id := c.Param("id")
	data, err := h.service.GetTopicById(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data")
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
		utils.Error(c, http.StatusInternalServerError, "failed to created data")
		return
	}

	utils.PrintJSON(data)

	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateTopic(c *gin.Context) {
	var req UpdateTopicRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	data, err := h.service.UpdateTopic(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to updated data")
		return
	}

	utils.Success(c, "success updated data", data)
}

func (h *handler) DeleteTopic(c *gin.Context) {
	var req TopicDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	data, err := h.service.DeleteTopic(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to deleted data")
		return
	}
	utils.Success(c, "success deleted data", data)
}
