package public

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

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
	sort := c.DefaultQuery("sort", "ASC")
	order := c.DefaultQuery("order", "id")
	search := c.DefaultQuery("search", "")
	topicParam := c.DefaultQuery("topics", "")
	isHighlightParam := c.DefaultQuery("is_highlight", "N")

	// Remove the square brackets from the parameter (e.g., "[1,2]" -> "1,2")
	topicParam = strings.Trim(topicParam, "[]")

	var topicIDs []int
	if topicParam == "" {
		topicIDs = []int{}
	} else {
		for _, topicStr := range strings.Split(topicParam, ",") {
			topicID, err := strconv.Atoi(topicStr)
			if err != nil {
				utils.Error(c, http.StatusBadRequest, "Invalid topic ID")
				return
			}
			topicIDs = append(topicIDs, topicID)
		}
	}

	// Call the GetAllPublicAuthors method of the service layer
	params := BlogPublicParams{
		Page:             page,
		Limit:            limit,
		Sort:             sort,
		Order:            order,
		Search:           search,
		Topics:           topicIDs,
		isHighlightParam: isHighlightParam,
	}

	// Validate the params using the binding tags
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	data, total_records, err := h.service.GetPublicBlogs(params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.PaginatedSuccess(c, "success get all data", data, page, limit, total_records)
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

func (h *handler) GetPublicTestimonials(c *gin.Context) {
	data, err := h.service.GetPublicTestimonials()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get all data", data)
}

func (h *handler) GetPublicTopics(c *gin.Context) {
	data, err := h.service.GetPublicTopics()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get all data", data)
}

func (h *handler) GetPublicProjects(c *gin.Context) {
	// Retrieve query parameters from the request
	page := utils.GetQueryParamInt(c, "page", 1)    // Default to page 1
	limit := utils.GetQueryParamInt(c, "limit", 10) // Default to 10 items per page
	sort := c.DefaultQuery("sort", "ASC")
	order := c.DefaultQuery("order", "id")
	search := c.DefaultQuery("search", "")
	isHighlightParam := c.DefaultQuery("is_highlight", "N")

	params := ProjectPublicParams{
		Page:             page,
		Limit:            limit,
		Sort:             sort,
		Order:            order,
		Search:           search,
		isHighlightParam: isHighlightParam,
	}

	// Validate the params using the binding tags
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	data, total_records, err := h.service.GetPublicProjects(params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.PaginatedSuccess(c, "success get all data", data, page, limit, total_records)
}

func (h *handler) GetPublicProjectBySlug(c *gin.Context) {
	slug := c.Param("slug")
	data, err := h.service.GetPublicProjectBySlug(slug)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get data", data)
}

func (h *handler) GetPublicTechnologies(c *gin.Context) {
	data, err := h.service.GetPublicTechnologies()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get all data", data)
}

func (h *handler) GetPublicAuthors(c *gin.Context) {
	data, err := h.service.GetPublicAuthors()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get all data", data)
}

func (h *handler) GetPublicExperiences(c *gin.Context) {
	data, err := h.service.GetPublicExperiences()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, "success get all data", data)
}

func (h *handler) UpdatePublicProjectStatistic(c *gin.Context) {
	var req ProjectStatisticUpdatePublicRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	data, err := h.service.UpdatePublicProjectStatistic(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success updated data", data)
}

func (h *handler) UpdatePublicBlogStatistic(c *gin.Context) {
	var req BlogStatisticUpdatePublicRequest

	if !utils.ValidateStruct(c, &req, c.ShouldBindJSON(&req)) {
		return
	}

	data, err := h.service.UpdatePublicBlogStatistic(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, "success updated data", data)
}
