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

	cleanupStaleMaterials(db)

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

// cleanupStaleMaterials 清除旧逻辑遗留的“已软删材质”：它们占用 name 唯一索引，
// 导致同名材质无法重建。仅清除无任何文物引用的记录。
func cleanupStaleMaterials(db *gorm.DB) {
	var staleIDs []uint
	if err := db.Unscoped().
		Model(&models.Material{}).
		Where("deleted_at IS NOT NULL").
		Pluck("id", &staleIDs).Error; err != nil {
		log.Printf("cleanup stale materials query failed: %v", err)
		return
	}
	for _, id := range staleIDs {
		var n int64
		if err := db.Model(&models.Find{}).Where("material_id = ?", id).Count(&n).Error; err != nil {
			log.Printf("cleanup count material %d failed: %v", id, err)
			continue
		}
		if n > 0 {
			continue
		}
		if err := db.Unscoped().Delete(&models.Material{}, id).Error; err != nil {
			log.Printf("cleanup hard delete material %d failed: %v", id, err)
		} else {
			log.Printf("cleaned up stale soft-deleted material id=%d", id)
		}
	}
}
