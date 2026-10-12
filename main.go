package main

import (
	"fmt"
	"log"
	"net/http"

	"simple-auction-system/internal/handlers"
	"simple-auction-system/internal/models"
	"simple-auction-system/internal/repositories"
	"simple-auction-system/internal/services"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 1. Kết nối Database SQLite
	db, err := gorm.Open(sqlite.Open("auction.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 2. Tự động Migrate bảng Users
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 3. Setup Dependency Injection
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// 4. Khởi tạo Router (Sử dụng ServeMux chuẩn của Go 1.22+)
	mux := http.NewServeMux()
	
	// Đăng ký các endpoints
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	// 5. Khởi chạy Server
	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}