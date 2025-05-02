package public

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

func (h *handler) GetAllPublicAuthors(c *gin.Context) {
	// Retrieve query parameters from the request
	page := utils.GetQueryParamInt(c, "page", 1)    // Default to page 1
	limit := utils.GetQueryParamInt(c, "limit", 10) // Default to 10 items per page
	sort := c.DefaultQuery("sort", "id")            // Default to sorting by "id"
	order := c.DefaultQuery("order", "ASC")         // Default to ascending order
	name := c.DefaultQuery("name", "")

	// Call the GetAllPublicAuthors method of the service layer
	params := AuthorPublicParams{
		Page:  page,
		Limit: limit,
		Sort:  sort,
		Order: order,
		Name:  name,
	}

	data, err := h.service.GetAllPublicAuthors(params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get all data", data)
}

func (h *handler) GetProfile(c *gin.Context) {
	data, err := h.service.GetProfile()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get data", data)
}

func (h *handler) GetPublicBlogs(c *gin.Context) {
	// Retrieve query parameters from the request
	page := utils.GetQueryParamInt(c, "page", 1)    // Default to page 1
	limit := utils.GetQueryParamInt(c, "limit", 10) // Default to 10 items per page
	sort := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "ASC")
	search := c.DefaultQuery("search", "")

	// Call the GetAllPublicAuthors method of the service layer
	params := BlogPublicParams{
		Page:   page,
		Limit:  limit,
		Sort:   sort,
		Order:  order,
		Search: search,
	}

	// Validate the params using the binding tags
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	data, err := h.service.GetPublicBlogs(params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get all data", data)
}

func (h *handler) GetPublicBlogBySlug(c *gin.Context) {
	slug := c.Param("slug")
	data, err := h.service.GetPublicBlogBySlug(slug)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get data", data)
}
