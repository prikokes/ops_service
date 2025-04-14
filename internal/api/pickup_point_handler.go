package api

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"ops_service/internal/model"
	"ops_service/internal/service"

	"github.com/gin-gonic/gin"
)

type PickupPointHandler struct {
	service *service.PickupPointService
}

func NewPickupPointHandler(service *service.PickupPointService) *PickupPointHandler {
	return &PickupPointHandler{
		service: service,
	}
}

func (h *PickupPointHandler) Create(c *gin.Context) {
	var req CreatePickupPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	point, err := h.service.CreatePickupPoint(c.Request.Context(), req.City)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, model.ErrCityNotAllowed) {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, point)
}

func (h *PickupPointHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var idInt int
	if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	point, err := h.service.GetPickupPoint(c.Request.Context(), idInt)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "pickup point not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, point)
}

func (h *PickupPointHandler) List(c *gin.Context) {
	var req GetPickupPointsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	points, total, err := h.service.ListPickupPoints(
		c.Request.Context(),
		req.Page,
		req.PageSize,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	response := PaginatedResponse{
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
		Data:       points,
	}

	c.JSON(http.StatusOK, response)
}
