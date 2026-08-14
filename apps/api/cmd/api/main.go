package main

import (
	"log"

	"github.com/your-org/splitr/apps/api/internal/config"
	"github.com/your-org/splitr/apps/api/internal/database"
	httpapi "github.com/your-org/splitr/apps/api/internal/http"
	"github.com/your-org/splitr/apps/api/internal/repository"
	"github.com/your-org/splitr/apps/api/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}

	store := repository.NewStore(db)
	auth := service.NewAuthService(store, cfg.SessionTTL, cfg.PasswordPepper)
	groups := service.NewGroupService(store)

	app := httpapi.NewServer(cfg, auth, groups)
	if err := app.Listen(cfg.APIAddr); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
