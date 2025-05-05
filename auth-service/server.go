package main

import (
	"auth-service/internal/delivery/http"
	"auth-service/internal/domain"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/pkg/config"

	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	db.Exec("PRAGMA journal_mode = WAL;")

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.RefreshToken{},
	); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}

	userRepo := repository.NewUserRepo(db)
	tokenRepo := repository.NewTokenRepo(db)

	authSvc := service.NewAuthService(
		userRepo, tokenRepo,
		cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL,
	)

	authH := http.NewAuthHandler(authSvc)

	router := gin.Default()
	router.POST("/api/v1/auth/register", authH.Register)
	router.POST("/api/v1/auth/login", authH.Login)
	router.POST("/api/v1/auth/refresh", authH.Refresh)
	router.POST("/api/v1/auth/logout", authH.Logout)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server run failed: %v", err)
	}
}
