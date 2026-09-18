package main

import (
	"log"
	"time"

	"digcatalog/internal/config"
	"digcatalog/internal/handlers"
	"digcatalog/internal/middleware"
	"digcatalog/internal/models"
	"digcatalog/internal/seed"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := config.Load()

	var db *gorm.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			sqlDB, e := db.DB()
			if e == nil && sqlDB.Ping() == nil {
				break
			}
			err = e
		}
		log.Printf("waiting for database... (%d/30): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Site{},
		&models.Unit{},
		&models.Material{},
		&models.Find{},
	); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	seed.Run(db)

	h := handlers.New(db, cfg.JWTSecret)
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/auth/login", h.Login)

		auth := api.Group("")
		auth.Use(middleware.AuthRequired(cfg.JWTSecret))
		{
			auth.GET("/auth/me", h.Me)
			auth.GET("/overview", h.Overview)

			auth.GET("/sites", h.ListSites)
			auth.GET("/sites/:id", h.GetSite)
			auth.POST("/sites", h.CreateSite)
			auth.PUT("/sites/:id", h.UpdateSite)
			auth.DELETE("/sites/:id", h.DeleteSite)

			auth.GET("/units", h.ListUnits)
			auth.GET("/units/:id", h.GetUnit)
			auth.POST("/units", h.CreateUnit)
			auth.PUT("/units/:id", h.UpdateUnit)
			auth.DELETE("/units/:id", h.DeleteUnit)

			auth.GET("/materials", h.ListMaterials)
			auth.POST("/materials", h.CreateMaterial)
			auth.PUT("/materials/:id", h.UpdateMaterial)
			auth.DELETE("/materials/:id", h.DeleteMaterial)

			auth.GET("/finds", h.ListFinds)
			auth.GET("/finds/:id", h.GetFind)
			auth.POST("/finds", h.CreateFind)
			auth.PUT("/finds/:id", h.UpdateFind)
			auth.DELETE("/finds/:id", h.DeleteFind)
		}
	}

	log.Printf("DigCatalog backend listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
