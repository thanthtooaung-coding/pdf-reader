package app

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/config"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/database"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/handler"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/otp"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/repository"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/routes"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/service"
	"github.com/thanthtooaung-coding/pdf-reader/backend/pkg/email"
	"github.com/thanthtooaung-coding/pdf-reader/backend/pkg/storage"
)

type Application struct {
	App         *fiber.App
	Config      *config.Config
	DB          *gorm.DB
	Redis       *redis.Client
	StoragePath string
}

func Bootstrap(cfg *config.Config) (*Application, error) {
	log := newLogger(cfg.LogLevel)

	db, err := database.Connect(cfg, log)
	if err != nil {
		return nil, err
	}

	if cfg.SeedDefaultAdmin {
		if err := database.SeedDefaultAdmin(db, log, cfg.DefaultAdminEmail, cfg.DefaultAdminPassword); err != nil {
			log.WithError(err).Error("seed default admin (continuing)")
		}
	}

	redisClient, err := otp.NewRedisClient(
		cfg.RedisURL, cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, cfg.RedisDB,
	)
	if err != nil {
		return nil, err
	}

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		redisClient.Close()
		return nil, err
	}

	store, err := storage.NewLocalStorage(cfg.StoragePath)
	if err != nil {
		redisClient.Close()
		return nil, err
	}

	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	workspaceRepo := repository.NewWorkspaceRepository(db)
	fileRepo := repository.NewFileRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	aiJobRepo := repository.NewAIJobRepository(db)

	otpStore := otp.NewStore(redisClient, cfg.JWTSecret)
	mailer := email.NewResendSender(email.ResendConfig{
		APIKey:  cfg.ResendAPIKey,
		BaseURL: cfg.ResendBaseURL,
		From:    cfg.EmailFrom,
	})

	authSvc := service.NewAuthService(log, cfg, db, userRepo, roleRepo, otpStore, mailer)
	userSvc := service.NewUserService(log, userRepo)
	workspaceSvc := service.NewWorkspaceService(log, workspaceRepo, fileRepo)
	fileSvc := service.NewFileService(log, fileRepo, workspaceRepo, store)
	commentSvc := service.NewCommentService(log, commentRepo, fileRepo)
	aiJobSvc := service.NewAIJobService(log, aiJobRepo, fileRepo, workspaceRepo)

	handlers := routes.Handlers{
		Auth:      handler.NewAuthHandler(authSvc),
		User:      handler.NewUserHandler(userSvc),
		Workspace: handler.NewWorkspaceHandler(workspaceSvc),
		File:      handler.NewFileHandler(fileSvc, cfg.MaxUploadMB),
		Comment:   handler.NewCommentHandler(commentSvc),
		AIJob:     handler.NewAIJobHandler(aiJobSvc),
	}

	fiberApp := fiber.New(fiber.Config{
		AppName:      "pdf-reader-backend",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
		BodyLimit:    cfg.MaxUploadMB * 1024 * 1024,
	})

	fiberApp.Use(recover.New())
	fiberApp.Use(fiberlogger.New())
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet, fiber.MethodPost, fiber.MethodPut,
			fiber.MethodPatch, fiber.MethodDelete, fiber.MethodOptions,
		}, ","),
	}))

	fiberApp.Static("/storage", store.BasePath())
	routes.Register(fiberApp, cfg, handlers)

	return &Application{
		App:         fiberApp,
		Config:      cfg,
		DB:          db,
		Redis:       redisClient,
		StoragePath: cfg.StoragePath,
	}, nil
}

func (a *Application) Close() error {
	if a.Redis != nil {
		return a.Redis.Close()
	}
	return nil
}

func newLogger(level string) *logrus.Logger {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.ErrorLevel
	}
	l.SetLevel(lvl)
	return l
}
