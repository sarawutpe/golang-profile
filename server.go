package main

import (
	"main/api"
	"main/db"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

// handling middleware
func authorized(c *gin.Context) {
	token := c.Query("token")
	if token != "7FC2D72AB1D9E" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "access denied!"})
		c.Abort()
	}
	c.Next()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	corsApi := router.Group("/api", authorized)
	{
		corsApi.GET("/avatar/:id", api.GetAvatarById)
		corsApi.PUT("/avatar/:id", api.UpdateAvatar)
	}
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))
	router.Static("/public", "./public")
	db.SetupDB()
	// start port
	port := os.Getenv("PORT")
	if port == "" {
		router.Run(":3000")
	} else {
		router.Run(":" + port)
	}

}
