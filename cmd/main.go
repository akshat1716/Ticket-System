package main

import (
	"log"

	"ticket-system/config"
	"ticket-system/controllers"
	"ticket-system/database"
	"ticket-system/middleware"
	"ticket-system/routes"
	"ticket-system/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	authService := services.NewAuthService(db, cfg.JWTSecret)
	ticketService := services.NewTicketService(db)

	authController := controllers.NewAuthController(authService)
	ticketController := controllers.NewTicketController(ticketService)

	router := gin.New()
	router.Use(middleware.Logging())
	router.Use(middleware.Recovery())

	// Serve static frontend files
router.Static("/static", "./frontend")

// Serve homepage
router.GET("/", func(c *gin.Context) {
	c.File("./frontend/index.html")
})

	routes.Setup(router, authController, ticketController, cfg.JWTSecret)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
