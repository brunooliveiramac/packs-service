package httpapi

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/brunooliveiramac/packs-service/internal/pack"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/brunooliveiramac/packs-service/internal/platform/dataprovider/database"
)

func RegisterRoutes(router *gin.Engine, sizesRepo *database.SizesRepository) {
	router.Use(gin.Logger())
	// Allow all origins for simplicity during deployment; tighten later if needed.
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	// calculation endpoint
	router.POST("/api/packs/calc", func(c *gin.Context) {
		var req struct {
			Quantity int `json:"quantity"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		// Always rely on sizes enrolled in the database
		sizes, err := sizesRepo.ActiveSizes(context.Background())
		if err != nil || len(sizes) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no pack sizes configured"})
			return
		}
		calc, err := pack.NewCalculator(sizes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		breakdown, shipped, err := calc.Calculate(req.Quantity)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// convert breakdown to array for stable JSON
		type item struct{ Size int `json:"size"`; Count int `json:"count"` }
		resp := struct {
			Requested int    `json:"requested"`
			Shipped   int    `json:"shipped"`
			Packs     []item `json:"packs"`
		}{Requested: req.Quantity, Shipped: shipped}
		for sz, ct := range breakdown {
			resp.Packs = append(resp.Packs, item{Size: sz, Count: ct})
		}
		c.JSON(http.StatusOK, resp)
	})

	// admin endpoints for managing pack sizes
	router.GET("/api/pack-sizes", func(c *gin.Context) {
		rows, err := sizesRepo.List(c, c.Query("all") == "true")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rows)
	})
	router.POST("/api/pack-sizes", func(c *gin.Context) {
		var req struct {
			Size   int  `json:"size"`
			Active *bool `json:"active"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.Size <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		active := true
		if req.Active != nil { active = *req.Active }
		if err := sizesRepo.Upsert(c, req.Size, active); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	})
	router.PATCH("/api/pack-sizes/:size", func(c *gin.Context) {
		var req struct{ Active bool `json:"active"` }
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		size, err := strconv.Atoi(c.Param("size"))
		if err != nil || size <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid size"})
			return
		}
		if err := sizesRepo.SetActive(c, size, req.Active); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})
	router.DELETE("/api/pack-sizes/:size", func(c *gin.Context) {
		size, err := strconv.Atoi(c.Param("size"))
		if err != nil || size <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid size"})
			return
		}
		if err := sizesRepo.Delete(c, size); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})
}


