package routes

import (
	"ghostline-backend/internal/handlers"
	"ghostline-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func Register(
	app *fiber.App,
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	adminHandler *handlers.AdminHandler,
	callHandler *handlers.CallHandler,
	userHandler *handlers.UserHandler,
	postHandler *handlers.PostHandler,
	likeHandler *handlers.LikeHandler,
	chatHandler *handlers.ChatHandler,
	websocketHandler *handlers.WebSocketHandler,
	jwtMiddleware *middleware.JWTMiddleware,
	adminMiddleware *middleware.AdminMiddleware,
	rateLimiter *middleware.RateLimiter,
) {
	app.Get("/health", healthHandler.Live)

	api := app.Group("/api")
	registerAPIRoutes(api, healthHandler, authHandler, adminHandler, callHandler, userHandler, postHandler, likeHandler, chatHandler, jwtMiddleware, adminMiddleware, rateLimiter)
	registerWebSocketRoutes(app, websocketHandler)
}
