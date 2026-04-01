package main

import (
	httpdelivery "delivery-service/internal/delivery/http"
	"delivery-service/internal/repository/postgres"
	"delivery-service/internal/usecase"
	pgconn "delivery-service/pkg/postgres"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := pgconn.New(pgconn.Config{
		Host:     getenv("DB_HOST", "localhost"),
		Port:     getenv("DB_PORT", "5432"),
		User:     getenv("DB_USER", "postgres"),
		Password: getenv("DB_PASSWORD", "postgres"),
		DBName:   getenv("DB_NAME", "delivery"),
		SSLMode:  getenv("DB_SSLMODE", "disable"),
	})
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer db.Close()

	orderRepo := postgres.NewOrderRepo(db)
	courierRepo := postgres.NewCourierRepo(db)

	orderUC := usecase.NewOrderUseCase(orderRepo, courierRepo)
	courierUC := usecase.NewCourierUseCase(courierRepo)

	orderH := httpdelivery.NewOrderHandler(orderUC)
	courierH := httpdelivery.NewCourierHandler(courierUC)

	router := httpdelivery.NewRouter(orderH, courierH)

	addr := getenv("ADDR", ":8080")
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
