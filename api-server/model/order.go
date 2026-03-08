package model

import "time"

// Order represents an order in the system.
type Order struct {
	ID                int64     `json:"id"`
	ProductName       string    `json:"product_name"`
	Quantity          int       `json:"quantity"`
	Note              *string   `json:"note"`
	Status            string    `json:"status"`
	ConfirmationToken string    `json:"confirmation_token,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// CreateOrderRequest represents the request body for creating an order.
type CreateOrderRequest struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Note        *string `json:"note"`
}

// UpdateOrderRequest represents the request body for updating an order.
type UpdateOrderRequest struct {
	Quantity *int    `json:"quantity"`
	Note     *string `json:"note"`
}

// ConfirmOrderRequest represents the request body for confirming an order.
type ConfirmOrderRequest struct {
	ConfirmationToken string `json:"confirmation_token"`
}

// ConfirmOrderResponse represents the response for confirming an order.
type ConfirmOrderResponse struct {
	ID          int64     `json:"id"`
	Status      string    `json:"status"`
	ConfirmedAt time.Time `json:"confirmed_at"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"`
}
