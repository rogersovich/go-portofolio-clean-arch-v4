package reading_time

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

func (h *handler) GetAll(c *gin.Context) {
	data, err := h.service.GetAllReadingTimes()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data", err)
		return
	}
	utils.Success(c, "success get all data", data)
}

func (h *handler) GetReadingTimeById(c *gin.Context) {
	// htmlInput1 := "<p><b>Hello world!</b> This is a simple test.</p><p>Another paragraph.</p>"
	// stats1 := utils.ExtractHTMLtoStatistics(htmlInput1)
	// fmt.Printf("Input 1:\n%s\n", htmlInput1)
	// fmt.Printf("Stats 1: %+v\n\n", stats1)

	id := c.Param("id")
	data, err := h.service.GetReadingTimeById(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to get data", err)
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

	data, err := h.service.CreateReadingTime(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to created data", err)
		return
	}

	utils.PrintJSON(data)

	utils.Success(c, "success get data", data)
}

func (h *handler) UpdateReadingTime(c *gin.Context) {
	var req UpdateReadingTimeRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	data, err := h.service.UpdateReadingTime(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to updated data", err)
		return
	}

	utils.Success(c, "success updated data", data)
}

func (h *handler) DeleteReadingTime(c *gin.Context) {
	var req ReadingTimeDeleteRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	id := req.ID

	data, err := h.service.DeleteReadingTime(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to deleted data", err)
		return
	}
	utils.Success(c, "success deleted data", data)
}
