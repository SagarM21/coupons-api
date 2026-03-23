package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"e-commerce/api/routes"
	"e-commerce/pkg/config"
	"e-commerce/pkg/db"
	"e-commerce/pkg/services"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.Init("config.yml"); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Create keyspace if it doesn't exist
	log.Println("creating keyspace if not exists...")
	for i := 0; i < 30; i++ {
		err := db.CreateKeyspace(config.App.Database)
		if err == nil {
			break
		}
		log.Printf("waiting for scylla to be ready... (%d/30): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	database, err := db.Init()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := database.CreateTables(); err != nil {
		log.Fatalf("failed to create tables: %v", err)
	}

	couponService := &services.CouponService{DB: database}
	couponHandler := &routes.CouponHandler{Service: couponService, DB: database}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	routes.RegisterCouponRoutes(r, couponHandler)

	addr := fmt.Sprintf("%s:%s", config.App.Api.Host, config.App.Api.Port)
	log.Printf("server starting on %s", addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
