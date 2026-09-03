package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"mariadiezmaback/internal/config"
	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/handler"
	"mariadiezmaback/internal/mailer"
	"mariadiezmaback/internal/middleware"
	"mariadiezmaback/internal/repository"
	"mariadiezmaback/internal/repository/memory"
	"mariadiezmaback/internal/repository/postgres"
	"mariadiezmaback/internal/service"
)

func main() {
	cfg := config.Load()

	// Configure structured logger
	var logHandler slog.Handler
	if cfg.Env == "production" {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	logger.Info("Starting MariaDiezmaBack API server",
		slog.String("env", cfg.Env),
		slog.String("port", cfg.Port),
		slog.String("db_driver", cfg.DBDriver),
	)

	// Initialize repositories based on configured driver
	var (
		userRepo  repository.UserRepository
		reqRepo   repository.RequestRepository
		colRepo   repository.CollectionRepository
		dressRepo repository.DressRepository
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if cfg.DBDriver == "postgres" {
		logger.Info("Connecting to PostgreSQL database...", slog.String("database_url", cfg.DatabaseURL))
		pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
		if err != nil {
			logger.Error("FATAL: Failed to connect to PostgreSQL. Make sure PostgreSQL is running (e.g. 'make docker-up' or local service) and DATABASE_URL is valid.", slog.Any("error", err))
			os.Exit(1)
		}
		defer pool.Close()
		logger.Info("✓ Connected to PostgreSQL successfully - Serving live data from database")
		userRepo = postgres.NewUserRepository(pool)
		reqRepo = postgres.NewRequestRepository(pool)
		colRepo = postgres.NewCollectionRepository(pool)
		dressRepo = postgres.NewDressRepository(pool)
	} else {
		logger.Info("Running with in-memory repository (ideal for testing and rapid local dev)")
		userRepo = memory.NewUserRepository()
		reqRepo = memory.NewRequestRepository()
		colRepo = memory.NewCollectionRepository()
		dressRepo = memory.NewDressRepository()
	}

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpirationHours)
	reqService := service.NewRequestService(reqRepo)
	mailService := mailer.New(cfg, logger)
	apptService := service.NewAppointmentService(reqRepo, mailService, cfg.NotificationEmail, logger)
	colService := service.NewCollectionService(colRepo)
	dressService := service.NewDressService(dressRepo)

	// Ensure default administrator account exists
	if err := authService.EnsureAdminUser(ctx, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		logger.Error("Failed to seed admin user", slog.Any("error", err))
	} else {
		logger.Info("Admin user ready", slog.String("email", cfg.AdminEmail))
	}

	// Seed default collections if empty
	if err := colService.EnsureDefaultCollections(ctx); err != nil {
		logger.Error("Failed to seed default collections", slog.Any("error", err))
	} else {
		logger.Info("Collections initialized")
	}

	// Seed default dresses if empty
	if err := dressService.EnsureDefaultDresses(ctx); err != nil {
		logger.Error("Failed to seed default dresses", slog.Any("error", err))
	} else {
		logger.Info("Dresses initialized")
	}

	// Handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)
	reqHandler := handler.NewRequestHandler(reqService)
	apptHandler := handler.NewAppointmentHandler(apptService)
	colHandler := handler.NewCollectionHandler(colService)
	dressHandler := handler.NewDressHandler(dressService)

	// Router setup
	r := chi.NewRouter()

	// Global middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.CORS(cfg.AllowedOrigins))

	// Health check
	r.Get("/health", healthHandler.Health)

	// API Routes
	r.Route("/api/v1", func(api chi.Router) {
		// Public Auth
		api.Post("/auth/login", authHandler.Login)

		// Public submission of general requests
		api.Post("/requests", reqHandler.Create)

		// Public submission of appointments (with email notification)
		api.Post("/appointments", apptHandler.Create)
		api.Post("/citas", apptHandler.Create) // Alias en español

		// Public collections for frontend web
		api.Get("/collections", colHandler.List)
		api.Get("/colecciones", colHandler.List) // Alias en español

		// Public dresses for frontend web
		api.Get("/dresses", dressHandler.List)
		api.Get("/vestidos", dressHandler.List) // Alias en español
		api.Get("/dresses/detail", dressHandler.GetDetail)
		api.Get("/vestidos/detalle", dressHandler.GetDetail) // Alias en español

		// Protected Backoffice Routes
		api.Group(func(backoffice chi.Router) {
			backoffice.Use(middleware.Auth(authService))

			// Current user info
			backoffice.Get("/auth/me", authHandler.Me)

			// Request Management
			backoffice.Get("/requests", reqHandler.List)
			backoffice.Get("/requests/{id}", reqHandler.GetByID)
			backoffice.Patch("/requests/{id}/status", reqHandler.UpdateStatus)
			backoffice.Put("/requests/{id}", reqHandler.Update)

			// Admin/Manager only: Delete requests
			backoffice.With(middleware.RequireRoles(domain.RoleAdmin, domain.RoleManager)).
				Delete("/requests/{id}", reqHandler.Delete)
		})
	})

	// Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Server shutdown channel
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		logger.Info("Shutting down server...", slog.String("signal", s.String()))

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer shutdownCancel()

		shutdownError <- server.Shutdown(shutdownCtx)
	}()

	logger.Info("Server listening", slog.String("addr", server.Addr))
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Server error", slog.Any("error", err))
		os.Exit(1)
	}

	if err := <-shutdownError; err != nil {
		logger.Error("Graceful shutdown error", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("Server stopped gracefully")
}
