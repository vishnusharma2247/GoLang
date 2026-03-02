package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"http-learning/internal/config"
	"http-learning/internal/repository/postgres"
	"http-learning/internal/service"
	httptransport "http-learning/internal/transport/http"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	dbPool, err := postgres.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer dbPool.Close()

	userRepo := postgres.NewUserRepository(dbPool)
	authService := service.NewAuthService(userRepo)
	authHandler := httptransport.NewAuthHandler(authService)

	mux := httptransport.NewMux(authHandler)

	server := &http.Server{
		Addr:    cfg.AppAddr,
		Handler: mux,
	}

	log.Printf("starting server on %s", cfg.AppAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
