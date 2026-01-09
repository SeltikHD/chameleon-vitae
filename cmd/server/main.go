// Package main is the entrypoint for the Chameleon Vitae server.
//
//	@title			Chameleon Vitae API
//	@version		1.0.0
//	@description	AI-powered resume engineering using Hexagonal Architecture. Tailors resumes to specific job descriptions using LLM technology.
//
//	@contact.name	Chameleon Vitae Team
//	@contact.url	https://github.com/SeltikHD/chameleon-vitae
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@host			localhost:8080
//	@BasePath		/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Firebase JWT token. Format: "Bearer {token}"
//
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	// Import generated swagger docs
	_ "github.com/SeltikHD/chameleon-vitae/docs"

	// Adapters
	httpAdapter "github.com/SeltikHD/chameleon-vitae/internal/adapters/primary/http"
	"github.com/SeltikHD/chameleon-vitae/internal/adapters/secondary/firebase"
	"github.com/SeltikHD/chameleon-vitae/internal/adapters/secondary/gotenberg"
	"github.com/SeltikHD/chameleon-vitae/internal/adapters/secondary/groq"
	"github.com/SeltikHD/chameleon-vitae/internal/adapters/secondary/jina"
	"github.com/SeltikHD/chameleon-vitae/internal/adapters/secondary/postgres"
	"github.com/SeltikHD/chameleon-vitae/internal/adapters/secondary/storage"

	// Config and Services
	"github.com/SeltikHD/chameleon-vitae/internal/config"
	"github.com/SeltikHD/chameleon-vitae/internal/core/services"
)

func main() {
	// Initialize context for startup operations
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize structured logger with zerolog
	initLogger(cfg)

	log.Info().
		Str("version", cfg.App.Version).
		Str("environment", cfg.App.Environment).
		Msg("Starting Chameleon Vitae server")

	// Initialize adapters
	adapters, err := initializeAdapters(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize adapters")
	}
	defer adapters.Close()

	// Initialize services
	svc := initializeServices(adapters)

	// Initialize HTTP router
	routerCfg := httpAdapter.RouterConfig{
		EnableSwagger:   cfg.Server.EnableSwagger,
		EnableProfiling: cfg.Server.EnableProfiling,
		RequestTimeout:  cfg.Server.WriteTimeout,
		MaxRequestSize:  cfg.Server.MaxRequestSize,
		AllowedOrigins:  cfg.Server.AllowedOrigins,
		BaseURL:         fmt.Sprintf("http://%s:%d", cfg.Server.Host, cfg.Server.Port),
	}

	router := httpAdapter.NewRouter(routerCfg, httpAdapter.Services{
		UserService:       svc.User,
		ExperienceService: svc.Experience,
		BulletService:     svc.Bullet,
		SkillService:      svc.Skill,
		ResumeService:     svc.Resume,
	})

	// Set up authentication middleware
	router.SetAuthMiddleware(adapters.Firebase, adapters.DB.UserRepository())

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Info().
			Str("host", cfg.Server.Host).
			Int("port", cfg.Server.Port).
			Bool("swagger", cfg.Server.EnableSwagger).
			Msg("HTTP server listening")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server stopped gracefully")
}

// initLogger initializes the zerolog logger based on configuration.
func initLogger(cfg *config.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Set log level
	if cfg.App.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Use pretty console output for development, JSON for production
	if cfg.IsDevelopment() {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}
}

// Adapters holds all initialized adapters.
type Adapters struct {
	DB        *postgres.DB
	Firebase  *firebase.Adapter
	Groq      *groq.Client
	Gotenberg *gotenberg.Client
	Jina      *jina.Client
	Storage   *storage.LocalStorage
}

// Close closes all adapters gracefully.
func (a *Adapters) Close() {
	if a.DB != nil {
		a.DB.Close()
		log.Debug().Msg("Database connection closed")
	}
	if a.Firebase != nil {
		if err := a.Firebase.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close Firebase adapter")
		}
	}
	if a.Groq != nil {
		if err := a.Groq.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close Groq client")
		}
	}
	if a.Gotenberg != nil {
		if err := a.Gotenberg.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close Gotenberg client")
		}
	}
	if a.Jina != nil {
		if err := a.Jina.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close Jina client")
		}
	}
	if a.Storage != nil {
		if err := a.Storage.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close storage")
		}
	}
	log.Info().Msg("All adapters closed")
}

// initializeAdapters initializes all secondary adapters.
func initializeAdapters(ctx context.Context, cfg *config.Config) (*Adapters, error) {
	adapters := &Adapters{}

	// Initialize PostgreSQL
	log.Info().Msg("Connecting to PostgreSQL...")
	dbCfg := postgres.Config{
		Host:              cfg.Database.Host,
		Port:              cfg.Database.Port,
		User:              cfg.Database.User,
		Password:          cfg.Database.Password, // pragma: allowlist secret
		Database:          cfg.Database.Database,
		SSLMode:           cfg.Database.SSLMode,
		MaxConns:          cfg.Database.MaxOpenConns,
		MinConns:          cfg.Database.MaxIdleConns,
		MaxConnLifetime:   cfg.Database.ConnMaxLifetime,
		MaxConnIdleTime:   cfg.Database.ConnMaxIdleTime,
		HealthCheckPeriod: cfg.Database.HealthCheckPeriod,
	}
	db, err := postgres.New(ctx, dbCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	adapters.DB = db
	log.Info().Msg("PostgreSQL connected successfully")

	// Initialize Firebase
	log.Info().Msg("Initializing Firebase authentication...")
	fbCfg := firebase.Config{
		ProjectID:       cfg.Firebase.ProjectID,
		CredentialsFile: cfg.Firebase.CredentialsFile,
	}
	if cfg.Firebase.CredentialsJSON != "" {
		fbCfg.CredentialsJSON = []byte(cfg.Firebase.CredentialsJSON)
	}
	fb, err := firebase.New(ctx, fbCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase: %w", err)
	}
	adapters.Firebase = fb
	log.Info().Msg("Firebase initialized successfully")

	// Initialize Groq
	log.Info().Msg("Initializing Groq AI provider...")
	groqCfg := groq.Config{
		APIKey:          cfg.Groq.APIKey, // pragma: allowlist secret
		ModelGeneration: cfg.Groq.DefaultModel,
		ModelAnalysis:   cfg.Groq.AnalysisModel,
		MaxRetries:      cfg.Groq.MaxRetries,
		Timeout:         cfg.Groq.RequestTimeout,
	}
	groqClient, err := groq.New(groqCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Groq: %w", err)
	}
	adapters.Groq = groqClient
	log.Info().Msg("Groq initialized successfully")

	// Initialize Gotenberg
	log.Info().Msg("Initializing Gotenberg PDF engine...")
	gotenCfg := gotenberg.Config{
		URL:     cfg.PDF.BaseURL,
		Timeout: cfg.PDF.Timeout,
	}
	gotenClient, err := gotenberg.New(gotenCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gotenberg: %w", err)
	}
	adapters.Gotenberg = gotenClient
	log.Info().Msg("Gotenberg initialized successfully")

	// Initialize Jina
	log.Info().Msg("Initializing Jina job parser...")
	jinaCfg := jina.Config{
		APIKey:  cfg.Jina.APIKey, // pragma: allowlist secret
		Timeout: cfg.Jina.Timeout,
	}
	jinaClient, err := jina.New(jinaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Jina: %w", err)
	}
	adapters.Jina = jinaClient
	log.Info().Msg("Jina initialized successfully")

	// Initialize Local Storage
	log.Info().Msg("Initializing file storage...")
	storageCfg := storage.LocalConfig{
		BasePath: cfg.Storage.LocalPath,
		BaseURL:  fmt.Sprintf("http://%s:%d/files", cfg.Server.Host, cfg.Server.Port),
	}
	localStorage, err := storage.NewLocalStorage(storageCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}
	adapters.Storage = localStorage
	log.Info().Msg("File storage initialized successfully")

	return adapters, nil
}

// Services holds all initialized services.
type Services struct {
	User       *services.UserService
	Experience *services.ExperienceService
	Bullet     *services.BulletService
	Skill      *services.SkillService
	Resume     *services.ResumeService
}

// initializeServices initializes all application services.
func initializeServices(adapters *Adapters) *Services {
	log.Info().Msg("Initializing services...")

	userService := services.NewUserService(
		adapters.DB.UserRepository(),
		adapters.Firebase,
	)

	experienceService := services.NewExperienceService(
		adapters.DB.ExperienceRepository(),
		adapters.DB.BulletRepository(),
	)

	bulletService := services.NewBulletService(
		adapters.DB.BulletRepository(),
		adapters.DB.ExperienceRepository(),
		adapters.Groq,
	)

	skillService := services.NewSkillService(
		adapters.DB.SkillRepository(),
		adapters.DB.SpokenLanguageRepository(),
	)

	resumeService := services.NewResumeService(
		adapters.DB.ResumeRepository(),
		adapters.DB.UserRepository(),
		adapters.DB.ExperienceRepository(),
		adapters.DB.BulletRepository(),
		adapters.DB.SkillRepository(),
		adapters.DB.SpokenLanguageRepository(),
		adapters.Groq,
		adapters.Gotenberg,
		adapters.Jina,
		adapters.Storage,
	)

	log.Info().Msg("All services initialized successfully")

	return &Services{
		User:       userService,
		Experience: experienceService,
		Bullet:     bulletService,
		Skill:      skillService,
		Resume:     resumeService,
	}
}
