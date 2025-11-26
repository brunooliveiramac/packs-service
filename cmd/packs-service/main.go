package main

import (
	"log"
	"os"

	"github.com/brunooliveiramac/packs-service/internal/platform/dataprovider/database"
	"github.com/brunooliveiramac/packs-service/internal/platform/metrics"
	httpapi "github.com/brunooliveiramac/packs-service/internal/platform/http"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"context"
)

func main() {
	router := gin.Default()

	metrics.MustRegister()
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// wire DB for sizes configuration
	ctx := context.Background()
	db, err := database.NewFromEnv(ctx)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	sizesRepo := database.NewSizesRepository(db.Pool())

	httpapi.RegisterRoutes(router, sizesRepo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	log.Printf("packs-service is running on :%s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}


