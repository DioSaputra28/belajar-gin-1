package main

import (
	"fmt"

	"github.com/DioSaputra28/belajar-gin-1/config"
	"github.com/DioSaputra28/belajar-gin-1/internal/addresses"
	"github.com/DioSaputra28/belajar-gin-1/internal/auth"
	"github.com/DioSaputra28/belajar-gin-1/internal/common/middleware"
	"github.com/DioSaputra28/belajar-gin-1/internal/contacts"
	"github.com/DioSaputra28/belajar-gin-1/internal/users"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	db := config.NewDB()
	if db == nil {
		panic("Failed to connect to database")
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Welcome to belajar-gin-1 API",
		})
	})

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

	contactRepo := contacts.NewContactRepository(db)
	contactSvc := contacts.NewContactService(contactRepo)
	contactHandler := contacts.NewContactHandler(contactSvc)

	contactAuth := router.Group("/contacts")
	contactAuth.Use(middleware.AuthMiddleware(authRepo))
	{
		contactAuth.GET("", contactHandler.GetContacts)
		contactAuth.POST("", contactHandler.CreateContact)
		contactAuth.PUT("/:id", contactHandler.UpdateContact)
		contactAuth.GET("/:id", contactHandler.FindContactById)
		contactAuth.DELETE("/:id", contactHandler.DeleteContact)
	}

	addressRepo := addresses.NewAddressRepository(db)
	addressSvc := addresses.NewAddressService(addressRepo)
	addressHandler := addresses.NewAddressHandler(addressSvc)

	addressAuth := router.Group("/addresses")
	addressAuth.Use(middleware.AuthMiddleware(authRepo))
	{
		addressAuth.GET("", addressHandler.GetAddresses)
		addressAuth.POST("", addressHandler.CreateAddress)
		addressAuth.PUT("/:id", addressHandler.UpdateAddress)
		addressAuth.GET("/:id", addressHandler.FindAddressById)
		addressAuth.DELETE("/:id", addressHandler.DeleteAddress)
	}

	fmt.Println("Starting server on :8081...")
	router.Run(":8081")
}
