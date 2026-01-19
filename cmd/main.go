package main

import (
	"github.com/DioSaputra28/belajar-gin-1/config"
	"github.com/DioSaputra28/belajar-gin-1/internal/auth"
	"github.com/DioSaputra28/belajar-gin-1/internal/common/middleware"
	"github.com/DioSaputra28/belajar-gin-1/internal/users"
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
	router.GET("/me", middleware.AuthMiddleware(authRepo), authHandler.Me)

	userRepo := users.NewUserRepository(db)
	userSvc := users.NewUserService(userRepo)
	userHandler := users.NewUserHandler(userSvc)

	userAuth := router.Group("/users")
	userAuth.Use(middleware.AuthMiddleware(authRepo))
	{
		userAuth.GET("", userHandler.GetUsers)
		userAuth.POST("", userHandler.CreateUser)
		userAuth.PUT("/:id", userHandler.UpdateUser)
		userAuth.GET("/:id", userHandler.FindUserById)
		userAuth.DELETE("/:id", userHandler.DeleteUser)
	}

	router.Run(":8081")
}