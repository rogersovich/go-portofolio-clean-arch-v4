package public

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

func (h *handler) GetAllPublicAuthors(c *gin.Context) {
	data, err := h.service.GetAllPublicAuthors()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get all data", data)
}
