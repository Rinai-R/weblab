package main

import (
	"log"
	"os"
	"weblab/internal/api"
	"weblab/internal/bootstrap"
	"weblab/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}

	container, err := bootstrap.BuildContainer(secret)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if container.DB != nil {
			_ = container.DB.Close()
		}
	}()

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS())

	api.RegisterRoutes(r, container)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
