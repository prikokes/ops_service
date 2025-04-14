package api

import (
	"ops_service/internal/model"
	"time"
)

type RegisterRequest struct {
	Email    string         `json:"email" binding:"required,email"`
	Password string         `json:"password" binding:"required,min=6"`
	Role     model.UserRole `json:"role" binding:"required,oneof=client moderator"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type CreatePickupPointRequest struct {
	City model.City `json:"city" binding:"required"`
}

type CreateReceiptRequest struct {
	PickupPointID int `json:"pickup_point_id" binding:"required"`
}

type AddProductRequest struct {
	Type model.ProductType `json:"type" binding:"required,oneof=электроника одежда обувь"`
}

type PaginationParams struct {
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1,max=100"`
}

type GetPickupPointsRequest struct {
	PaginationParams
	StartDate *time.Time `form:"start_date"`
	EndDate   *time.Time `form:"end_date"`
}

type PaginatedResponse struct {
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}
