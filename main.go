package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "simple-auction-system/docs"
	"simple-auction-system/internal/config"
	"simple-auction-system/internal/handlers"
	"simple-auction-system/internal/middleware"
	"simple-auction-system/internal/models"
	"simple-auction-system/internal/repositories"
	"simple-auction-system/internal/services"

	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Simple Auction System API
// @version 1.0
// @description REST API cho hệ thống đấu giá đơn giản
// @host localhost:8080
// @BasePath /api
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0o755); err != nil {
		log.Fatal("Failed to create data dir:", err)
	}

	// 1. Kết nối Database SQLite
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 2. Tự động Migrate bảng Users
	err = db.AutoMigrate(
		&models.User{},
		&models.Auction{},
		&models.Bid{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 3. Setup Dependency Injection
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	auctionRepository := repositories.NewAuctionRepository(db)
	auctionService := services.NewAuctionService(
		auctionRepository,
	)
	auctionHandler := handlers.NewAuctionHandler(
		auctionService,
	)

	// 4. Khởi tạo Router (Sử dụng ServeMux chuẩn của Go 1.22+)
	mux := http.NewServeMux()
	// docs
	mux.Handle(
		"/docs/",
		httpSwagger.Handler(
			httpSwagger.URL("/docs/doc.json"),
		),
	)

	// Đăng ký các endpoints
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.Handle(
		"GET /api/users/me",
		middleware.RequireAuth(authService)(
			http.HandlerFunc(authHandler.Me),
		),
	)

	mux.Handle(
		"POST /api/auctions",
		middleware.RequireAuth(authService)(
			http.HandlerFunc(auctionHandler.CreateAuction),
		),
	)
	mux.HandleFunc(
		"GET /api/auctions",
		auctionHandler.GetAllAuctions,
	)
	mux.HandleFunc(
		"GET /api/auctions/{id}",
		auctionHandler.GetAuctionByID,
	)

	// 5. Khởi chạy Server
	fmt.Println("Server is running on port " + cfg.Addr + " ...")
	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
