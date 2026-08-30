package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/your-org/splitr/apps/api/internal/config"
	"github.com/your-org/splitr/apps/api/internal/service"
)

type Server struct {
	cfg    config.Config
	auth   *service.AuthService
	groups *service.GroupService
}

func NewServer(cfg config.Config, auth *service.AuthService, groups *service.GroupService) *fiber.App {
	server := &Server{cfg: cfg, auth: auth, groups: groups}

	app := fiber.New(fiber.Config{
		AppName:      "splitr-api",
		ErrorHandler: server.errorHandler,
	})
	app.Use(logger.New(logger.Config{Format: `{"time":"${time}","status":${status},"latency":"${latency}","method":"${method}","path":"${path}"}\n`}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.WebOrigin,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET,POST,DELETE,OPTIONS",
	}))

	api := app.Group("/api/v1")
	api.Get("/health", server.health)
	api.Post("/auth/register", server.register)
	api.Post("/auth/login", server.login)
	api.Post("/auth/logout", server.requireUser, server.logout)
	api.Get("/me", server.requireUser, server.me)
	api.Get("/groups", server.requireUser, server.listGroups)
	api.Post("/groups", server.requireUser, server.createGroup)
	api.Get("/groups/:groupID/participants", server.requireUser, server.listParticipants)
	api.Post("/groups/:groupID/participants", server.requireUser, server.addParticipant)
	api.Get("/groups/:groupID/tags", server.requireUser, server.listTags)
	api.Post("/groups/:groupID/tags", server.requireUser, server.createTag)
	api.Get("/groups/:groupID/balances", server.requireUser, server.listBalances)
	api.Get("/groups/:groupID/expenses", server.requireUser, server.listExpenses)
	api.Post("/groups/:groupID/expenses", server.requireUser, server.createExpense)
	api.Get("/groups/:groupID/settlements", server.requireUser, server.listSettlements)
	api.Post("/groups/:groupID/settlements", server.requireUser, server.createSettlement)
	api.Delete("/groups/:groupID/settlements/:settlementID", server.requireUser, server.deleteSettlement)

	return app
}

func (s *Server) errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if fiberErr, ok := err.(*fiber.Error); ok {
		code = fiberErr.Code
	}
	return writeError(c, code, "INTERNAL_ERROR", "Something went wrong.", nil)
}

func (s *Server) health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}
