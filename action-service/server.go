package main

import (
	"action-service/internal/delivery/http"
	"action-service/internal/domain"
	"action-service/internal/repository"
	"action-service/internal/service"
	"action-service/pkg/config"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	db.Exec("PRAGMA journal_mode = WAL;")
	db.AutoMigrate(&domain.ActionCategory{}, &domain.DayAction{})

	catRepo := repository.NewCategoryRepo(db)
	actRepo := repository.NewActionRepo(db)
	catSvc := service.NewCategoryService(catRepo)
	actSvc := service.NewActionService(actRepo)
	catH := http.NewCategoryHandler(catSvc)
	actH := http.NewActionHandler(actSvc)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "time": time.Now()}) })

	api := router.Group("/api/v1")
	{
		cats := api.Group("/categories")
		{
			cats.GET("", catH.ListCategories)
			cats.POST("", catH.CreateCategory)
			// ...
		}
		acts := api.Group("/days/:date/actions")
		{
			acts.GET("", actH.ListActions)
			acts.POST("", actH.AddAction)
			// ...
		}
	}

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server run error: %v", err)
	}
}
