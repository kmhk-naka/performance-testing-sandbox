package handler

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kmhk-naka/performance-testing-sandbox/api-server/model"
	"github.com/kmhk-naka/performance-testing-sandbox/api-server/repository"
)

// OrderHandler handles HTTP requests for orders.
type OrderHandler struct {
	repo *repository.OrderRepository
	db   *sql.DB
}

// NewOrderHandler creates a new OrderHandler.
func NewOrderHandler(repo *repository.OrderRepository, db *sql.DB) *OrderHandler {
	return &OrderHandler{repo: repo, db: db}
}

// GetOrder handles GET /api/orders/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "注文IDが不正です")
		return
	}

	order, err := h.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "not_found", "注文が見つかりません")
			return
		}
		log.Printf("GetOrder error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "サーバーエラーが発生しました")
		return
	}

	writeJSON(w, http.StatusOK, order)
}

// CreateOrder handles POST /api/orders (stateless)
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req model.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_body", "リクエストボディが不正です")
		return
	}

	if req.ProductName == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "product_nameは必須です")
		return
	}
	if req.Quantity <= 0 {
		writeError(w, http.StatusBadRequest, "validation_error", "quantityは1以上を指定してください")
		return
	}

	token, err := generateToken()
	if err != nil {
		log.Printf("generateToken error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "サーバーエラーが発生しました")
		return
	}

	order := &model.Order{
		ProductName:       req.ProductName,
		Quantity:          req.Quantity,
		Note:              req.Note,
		Status:            "pending",
		ConfirmationToken: token,
	}

	if err := h.repo.Create(order); err != nil {
		log.Printf("CreateOrder error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "サーバーエラーが発生しました")
		return
	}

	// Don't expose confirmation_token in create response
	order.ConfirmationToken = ""
	writeJSON(w, http.StatusCreated, order)
}

// UpdateOrder handles PUT /api/orders/{id}
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "注文IDが不正です")
		return
	}

	// Check existence first
	_, err = h.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "not_found", "注文が見つかりません")
			return
		}
		log.Printf("UpdateOrder error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "サーバーエラーが発生しました")
		return
	}

	var req model.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_body", "リクエストボディが不正です")
		return
	}

	updated, err := h.repo.Update(id, &req)
	if err != nil {
		log.Printf("UpdateOrder error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "サーバーエラーが発生しました")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// ConfirmOrder handles POST /api/orders/{id}/confirm (stateful)
func (h *OrderHandler) ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "注文IDが不正です")
		return
	}

	// Check existence (condition B: order must exist)
	order, err := h.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "not_found", "注文が見つかりません")
			return
		}
		log.Printf("ConfirmOrder error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "サーバーエラーが発生しました")
		return
	}

	// Check if already confirmed
	if order.Status == "confirmed" {
		writeError(w, http.StatusConflict, "already_confirmed", "この注文は既に確定済みです")
		return
	}

	// Parse request body
	var req model.ConfirmOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_body", "リクエストボディが不正です")
		return
	}

	// Validate confirmation token (condition A: token must match)
	if req.ConfirmationToken == "" {
		writeError(w, http.StatusBadRequest, "missing_token", "confirmation_tokenは必須です")
		return
	}
	if req.ConfirmationToken != order.ConfirmationToken {
		writeError(w, http.StatusBadRequest, "invalid_token", "confirmation_tokenが不正です")
		return
	}

	// Confirm the order
	if err := h.repo.Confirm(id); err != nil {
		log.Printf("ConfirmOrder error: %v", err)
		writeError(w, http.StatusConflict, "confirm_failed", "注文の確定に失敗しました")
		return
	}

	// Fetch updated order to get confirmed_at
	updated, err := h.repo.GetByID(id)
	if err != nil {
		log.Printf("ConfirmOrder fetch error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "サーバーエラーが発生しました")
		return
	}

	resp := model.ConfirmOrderResponse{
		ID:          updated.ID,
		Status:      updated.Status,
		ConfirmedAt: updated.UpdatedAt,
	}
	writeJSON(w, http.StatusOK, resp)
}

// HealthCheck handles GET /health
func (h *OrderHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	dbStatus := "ok"
	if err := h.db.Ping(); err != nil {
		dbStatus = "error: " + err.Error()
	}

	resp := model.HealthResponse{
		Status: "ok",
		DB:     dbStatus,
	}
	writeJSON(w, http.StatusOK, resp)
}

// --- Helpers ---

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, errCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model.ErrorResponse{
		Error:   errCode,
		Message: message,
	})
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
