package main

import (
	"github.com/DioSaputra28/belajar-gin-1/config"
	"github.com/DioSaputra28/belajar-gin-1/internal/auth"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	db := config.NewDB()

	authRepo := auth.NewAuthRepository(db)
	authSvc := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authSvc)

	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/refresh-token", authHandler.RefreshToken)

	router.Run(":8080")
}