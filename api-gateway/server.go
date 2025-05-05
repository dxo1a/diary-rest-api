package main

import (
	"api-gateway/internal/middleware"
	"api-gateway/internal/routes"
	"api-gateway/pkg/config"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	router := gin.New()
	router.Use(
		middleware.Logging(),
	)

	// ---- OPEN AUTH ROUTES ----
	router.POST("/api/v1/auth/register", routes.ReverseProxy(cfg.AuthService))
	router.POST("/api/v1/auth/login", routes.ReverseProxy(cfg.AuthService))
	router.POST("/api/v1/auth/refresh", routes.ReverseProxy(cfg.AuthService))
	// ---- OPEN ROUTES ----

	router.Use(middleware.Auth(cfg.JWTSecret))

	// ---- RESTRICTED AUTH ROUTES ----
	router.POST("/api/v1/auth/logout", routes.ReverseProxy(cfg.AuthService))
	// ---- RESTRICTED AUTH ROUTES ----

	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "time": time.Now()}) })

	// регистрация всех прокси маршрутов под /api/v1
	routes.Register(router, cfg)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server run error: %v", err)
	}
}
