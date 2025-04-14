package api

import (
	"errors"
	"fmt"
	"net/http"
	"ops_service/internal/model"
	"ops_service/internal/service"

	"github.com/gin-gonic/gin"
)

type ReceiptHandler struct {
	service *service.ReceiptService
}

func NewReceiptHandler(service *service.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{
		service: service,
	}
}

func (h *ReceiptHandler) Create(c *gin.Context) {
	var req CreateReceiptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	receipt, err := h.service.CreateReceipt(c.Request.Context(), req.PickupPointID)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, model.ErrReceiptAlreadyClosed) {
			c.JSON(http.StatusConflict, gin.H{"error": "There is already an open receipt for this pickup point"})
			return
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, receipt)
}

func (h *ReceiptHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var idInt int
	if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	receipt, err := h.service.GetReceipt(c.Request.Context(), idInt)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "receipt not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, receipt)
}

func (h *ReceiptHandler) Close(c *gin.Context) {
	id := c.Param("id")
	var idInt int
	if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.CloseReceipt(c.Request.Context(), idInt); err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, model.ErrReceiptAlreadyClosed) {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Receipt closed successfully"})
}

func (h *ReceiptHandler) AddProduct(c *gin.Context) {
	var req AddProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pickupPointID, err := getPickupPointID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.AddProduct(c.Request.Context(), pickupPointID, req.Type)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, model.ErrNoOpenReceipt) {
			status = http.StatusBadRequest
			c.JSON(status, gin.H{"error": "No open receipt available. Please create a new receipt first."})
			return
		}

		if errors.Is(err, model.ErrReceiptAlreadyClosed) {
			status = http.StatusBadRequest
			c.JSON(status, gin.H{"error": "The receipt is already closed."})
			return
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ReceiptHandler) DeleteLastProduct(c *gin.Context) {
	pickupPointID, err := getPickupPointID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.DeleteLastProduct(c.Request.Context(), pickupPointID); err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, model.ErrNoOpenReceipt) {
			status = http.StatusBadRequest
			c.JSON(status, gin.H{"error": "No open receipt available."})
			return
		}

		if errors.Is(err, model.ErrReceiptAlreadyClosed) {
			status = http.StatusBadRequest
			c.JSON(status, gin.H{"error": "The receipt is already closed."})
			return
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Last product deleted successfully"})
}

func getPickupPointID(c *gin.Context) (int, error) {
	pickupPointID := c.Query("pickup_point_id")
	if pickupPointID == "" {
		return 0, errors.New("pickup_point_id is required")
	}

	var id int
	if _, err := fmt.Sscanf(pickupPointID, "%d", &id); err != nil {
		return 0, errors.New("invalid pickup_point_id")
	}

	return id, nil
}
