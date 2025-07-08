package main

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/cavidyrm/instawall/config"
	// --- New, Feature-Specific Imports ---
	categorydelivery "github.com/cavidyrm/instawall/internal/category/delivery/http"
	categoryRepo "github.com/cavidyrm/instawall/internal/category/repository/postgres"
	categoryUsecase "github.com/cavidyrm/instawall/internal/category/usecase"
	pagedelivery "github.com/cavidyrm/instawall/internal/page/delivery/http"
	pageRepo "github.com/cavidyrm/instawall/internal/page/repository/postgres"
	pageUsecase "github.com/cavidyrm/instawall/internal/page/usecase"
	// --- User Imports ---
	userdelivery "github.com/cavidyrm/instawall/internal/user/delivery/http"
	userRepo "github.com/cavidyrm/instawall/internal/user/repository/postgres"
	redisRepo "github.com/cavidyrm/instawall/internal/user/repository/redis"
	userUsecase "github.com/cavidyrm/instawall/internal/user/usecase"
	// --- Common Packages ---
	"github.com/cavidyrm/instawall/pkg/database"
	"github.com/cavidyrm/instawall/pkg/filestore"
	"github.com/cavidyrm/instawall/pkg/migration"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// 2. Initialize External Services
	db, err := database.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Fatalf("could not initialize postgres db: %v", err)
	}
	defer db.Close()

	migration.Run(db, cfg.Postgres.DBName)

	rdb := database.NewRedisClient(cfg.Redis)
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("could not connect to redis: %v", err)
	}

	fs, err := filestore.NewFileStore(cfg.MinIO)
	if err != nil {
		log.Fatalf("could not initialize minio filestore: %v", err)
	}

	// 3. Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 4. Initialize Repositories
	userRepository := userRepo.NewUserRepository(db)
	otpRepository := redisRepo.NewOTPRepository(rdb)
	pageRepository := pageRepo.NewPageRepository(db)
	categoryRepository := categoryRepo.NewCategoryRepository(db)

	// 5. Initialize Usecases
	userUC := userUsecase.NewUserUsecase(userRepository, otpRepository)
	pageUC := pageUsecase.NewPageUsecase(pageRepository, fs)
	categoryUC := categoryUsecase.NewCategoryUsecase(categoryRepository, fs)

	// 6. Register deliverys
	userdelivery.Registerdeliverys(e, userUC)
	pagedelivery.RegisterPagedeliverys(e, pageUC)
	categorydelivery.RegisterCategorydeliverys(e, categoryUC)

	// 7. Start Server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := e.Start(cfg.Server.Port); err != nil {
		e.Logger.Fatal(err)
	}
}
