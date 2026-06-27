package routes

import (
	"ticket-system/controllers"
	"ticket-system/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(
	router *gin.Engine,
	authController *controllers.AuthController,
	ticketController *controllers.TicketController,
	jwtSecret string,
) {
	router.GET("/health", controllers.Health)

	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	tickets := router.Group("/tickets")
	tickets.Use(middleware.JWTAuth(jwtSecret))
	{
		tickets.POST("", ticketController.Create)
		tickets.GET("", ticketController.List)
		tickets.GET("/:id", ticketController.GetByID)
		tickets.PATCH("/:id/status", ticketController.UpdateStatus)
	}

	router.Static("/static", "./frontend")
	router.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})
}
