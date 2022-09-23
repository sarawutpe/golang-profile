package main

import (
	"main/api"
	"main/db"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// handling middleware
func authorized(c *gin.Context) {
	c.Next()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.MaxMultipartMemory = 1 << 20
	corsApi := r.Group("/api", authorized)
	{
		corsApi.GET("/avatar/:id", api.GetAvatarById)
		corsApi.PUT("/avatar/:id", api.UpdateAvatar)
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT"},
		AllowHeaders:     []string{"Origin, X-Requested-With, Content-Type, Accept, Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	r.Static("/public", "./public")
	db.SetupDB()
	// start port
	port := os.Getenv("PORT")
	if port == "" {
		r.Run(":3000")
	} else {
		r.Run(":" + port)
	}
}
