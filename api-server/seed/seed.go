package seed

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
)

// SeedOrders inserts initial order data for load testing.
func SeedOrders(db *sql.DB, count int) error {
	// Check if data already exists
	var existingCount int
	err := db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&existingCount)
	if err != nil {
		return fmt.Errorf("failed to check existing orders: %w", err)
	}

	if existingCount >= count {
		log.Printf("Seed: %d orders already exist, skipping", existingCount)
		return nil
	}

	log.Printf("Seed: inserting %d orders...", count)

	products := []string{
		"ノートPC", "デスクトップPC", "モニター", "キーボード", "マウス",
		"ヘッドセット", "Webカメラ", "USBハブ", "SSD", "メモリ",
		"ルーター", "LANケーブル", "電源タップ", "デスクライト", "チェア",
		"デスク", "モニターアーム", "ケーブルトレー", "フットレスト", "リストレスト",
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		"INSERT INTO orders (product_name, quantity, note, status, confirmation_token) VALUES (?, ?, ?, 'pending', ?)",
	)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for i := 0; i < count; i++ {
		product := products[i%len(products)]
		quantity := (i % 10) + 1
		note := fmt.Sprintf("シードデータ #%d", i+1)

		token, err := generateSeedToken()
		if err != nil {
			return fmt.Errorf("failed to generate token: %w", err)
		}

		_, err = stmt.Exec(product, quantity, note, token)
		if err != nil {
			return fmt.Errorf("failed to insert order %d: %w", i+1, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Seed: inserted %d orders successfully", count)
	return nil
}

func generateSeedToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
