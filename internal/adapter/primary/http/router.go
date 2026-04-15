package http

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *URLHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/urls", handler.CreateURL)
		v1.GET("/urls/:code", handler.GetURL)
		v1.GET("/urls", handler.ListURLs)
	}

	router.GET("/r/:code", handler.Redirect)

	return router
}
