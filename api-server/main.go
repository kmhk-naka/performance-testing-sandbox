package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"

	"github.com/kmhk-naka/performance-testing-sandbox/api-server/handler"
	"github.com/kmhk-naka/performance-testing-sandbox/api-server/repository"
	"github.com/kmhk-naka/performance-testing-sandbox/api-server/seed"
)

func main() {
	// Database configuration
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "app")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "orders_db")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to database with retry
	var db *sql.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Waiting for database... (attempt %d/30)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Connected to database successfully")

	// Seed data
	if err := seed.SeedOrders(db, 100); err != nil {
		log.Printf("Warning: seed failed: %v", err)
	}

	// Setup repository and handler
	orderRepo := repository.NewOrderRepository(db)
	orderHandler := handler.NewOrderHandler(orderRepo, db)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	// Routes
	r.Get("/health", orderHandler.HealthCheck)

	r.Route("/api/orders", func(r chi.Router) {
		r.Post("/", orderHandler.CreateOrder)
		r.Get("/{id}", orderHandler.GetOrder)
		r.Put("/{id}", orderHandler.UpdateOrder)
		r.Post("/{id}/confirm", orderHandler.ConfirmOrder)
	})

	// Start server
	port := getEnv("PORT", "8080")
	addr := ":" + port
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
