package api

import (
	"catalog/app"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// NewApp initializes a new Fiber application with necessary settings, middleware, and routes.
func NewApp(app *app.App) *fiber.App {
	// Create a new Fiber instance
	router := fiber.New()

	// Use logger middleware to log HTTP requests
	router.Use(logger.New())

	// Define routes
	router.Get("/cars", getCars(app))
	router.Delete("/cars/:id", deleteCar(app))
	router.Put("/cars", updateCars(app))
	router.Post("/cars", addCars(app))

	return router
}
