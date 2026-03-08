package repository

import (
	"database/sql"
	"fmt"

	"github.com/kmhk-naka/performance-testing-sandbox/api-server/model"
)

// OrderRepository handles database operations for orders.
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new OrderRepository.
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// GetByID retrieves an order by its ID.
func (r *OrderRepository) GetByID(id int64) (*model.Order, error) {
	var order model.Order
	err := r.db.QueryRow(
		"SELECT id, product_name, quantity, note, status, confirmation_token, created_at, updated_at FROM orders WHERE id = ?",
		id,
	).Scan(
		&order.ID, &order.ProductName, &order.Quantity, &order.Note,
		&order.Status, &order.ConfirmationToken, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// Create inserts a new order into the database.
func (r *OrderRepository) Create(order *model.Order) error {
	result, err := r.db.Exec(
		"INSERT INTO orders (product_name, quantity, note, status, confirmation_token) VALUES (?, ?, ?, ?, ?)",
		order.ProductName, order.Quantity, order.Note, order.Status, order.ConfirmationToken,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	order.ID = id

	// Fetch the created_at and updated_at from DB
	return r.db.QueryRow(
		"SELECT created_at, updated_at FROM orders WHERE id = ?", id,
	).Scan(&order.CreatedAt, &order.UpdatedAt)
}

// Update updates an existing order.
func (r *OrderRepository) Update(id int64, req *model.UpdateOrderRequest) (*model.Order, error) {
	// Build dynamic update query
	setClauses := []string{}
	args := []interface{}{}

	if req.Quantity != nil {
		setClauses = append(setClauses, "quantity = ?")
		args = append(args, *req.Quantity)
	}
	if req.Note != nil {
		setClauses = append(setClauses, "note = ?")
		args = append(args, *req.Note)
	}

	if len(setClauses) == 0 {
		return r.GetByID(id)
	}

	query := "UPDATE orders SET "
	for i, clause := range setClauses {
		if i > 0 {
			query += ", "
		}
		query += clause
	}
	query += " WHERE id = ?"
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return r.GetByID(id)
}

// Confirm changes the order status to 'confirmed'.
func (r *OrderRepository) Confirm(id int64) error {
	result, err := r.db.Exec(
		"UPDATE orders SET status = 'confirmed' WHERE id = ? AND status = 'pending'",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to confirm order: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("order already confirmed or not found")
	}
	return nil
}
