// Package http provides the primary HTTP adapter using Chi router.
// This adapter implements the input port for the application, handling
// HTTP requests and translating them to service calls.
package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/SeltikHD/chameleon-vitae/internal/core/services"
)

// RouterConfig holds configuration for the HTTP router.
type RouterConfig struct {
	// EnableSwagger enables the Swagger UI endpoint.
	EnableSwagger bool

	// EnableProfiling enables pprof endpoints.
	EnableProfiling bool

	// RequestTimeout is the maximum duration for processing a request.
	RequestTimeout time.Duration

	// MaxRequestSize is the maximum size of request body in bytes.
	MaxRequestSize int64

	// AllowedOrigins for CORS configuration.
	AllowedOrigins []string

	// BaseURL is the base URL for the API (used in Swagger).
	BaseURL string
}

// DefaultRouterConfig returns sensible defaults for the router.
func DefaultRouterConfig() RouterConfig {
	return RouterConfig{
		EnableSwagger:   true,
		EnableProfiling: false,
		RequestTimeout:  60 * time.Second,
		MaxRequestSize:  10 * 1024 * 1024, // 10MB
		AllowedOrigins:  []string{"*"},
		BaseURL:         "http://localhost:8080",
	}
}

// Services holds all service dependencies for the HTTP handlers.
type Services struct {
	UserService       *services.UserService
	ExperienceService *services.ExperienceService
	BulletService     *services.BulletService
	SkillService      *services.SkillService
	ResumeService     *services.ResumeService
}

// Router wraps the Chi router and handlers.
type Router struct {
	mux      *chi.Mux
	config   RouterConfig
	services Services

	// Authentication middleware
	authMiddleware *authMiddleware

	// Handlers
	authHandler       *AuthHandler
	userHandler       *UserHandler
	experienceHandler *ExperienceHandler
	bulletHandler     *BulletHandler
	skillHandler      *SkillHandler
	languageHandler   *SpokenLanguageHandler
	resumeHandler     *ResumeHandler
	toolsHandler      *ToolsHandler
}

// NewRouter creates a new HTTP router with the given configuration and services.
func NewRouter(cfg RouterConfig, svc Services) *Router {
	r := &Router{
		mux:      chi.NewRouter(),
		config:   cfg,
		services: svc,
	}

	r.setupMiddleware()
	r.setupHandlers()
	r.setupRoutes()

	return r
}

// setupMiddleware configures the middleware stack.
func (r *Router) setupMiddleware() {
	// Request ID for tracing
	r.mux.Use(middleware.RequestID)

	// Real IP extraction (for proxied requests)
	r.mux.Use(middleware.RealIP)

	// Structured logging with zerolog
	r.mux.Use(ZerologLogger)

	// Panic recovery with pretty stack traces
	r.mux.Use(middleware.Recoverer)

	// Request timeout
	r.mux.Use(middleware.Timeout(r.config.RequestTimeout))

	// Request size limit
	r.mux.Use(middleware.RequestSize(r.config.MaxRequestSize))

	// Strip trailing slashes
	r.mux.Use(middleware.StripSlashes)

	// Heartbeat endpoint (before auth)
	r.mux.Use(middleware.Heartbeat("/ping"))

	// Content type enforcement for POST/PUT/PATCH
	r.mux.Use(ContentTypeJSON)

	// CORS configuration
	r.mux.Use(CORS(r.config.AllowedOrigins))
}

// setupHandlers initializes all HTTP handlers.
func (r *Router) setupHandlers() {
	r.authHandler = NewAuthHandler(r.services.UserService)
	r.userHandler = NewUserHandler(r.services.UserService)
	r.experienceHandler = NewExperienceHandler(r.services.ExperienceService)
	r.bulletHandler = NewBulletHandler(r.services.BulletService)
	r.skillHandler = NewSkillHandler(r.services.SkillService)
	r.languageHandler = NewSpokenLanguageHandler(r.services.SkillService) // Spoken languages are in SkillService
	r.resumeHandler = NewResumeHandler(r.services.ResumeService)
	r.toolsHandler = NewToolsHandler(r.services.ResumeService) // Tools use ResumeService for job parsing
}

// setupRoutes configures all API routes.
func (r *Router) setupRoutes() {
	// Health check (unauthenticated)
	r.mux.Get("/health", r.healthHandler)

	// Swagger documentation (if enabled)
	if r.config.EnableSwagger {
		swaggerHost := r.config.BaseURL
		if strings.Contains(swaggerHost, "0.0.0.0") {
			swaggerHost = strings.ReplaceAll(swaggerHost, "0.0.0.0", "localhost")
		}
		r.mux.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(swaggerHost+"/swagger/doc.json"),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("list"),
			httpSwagger.DomID("swagger-ui"),
			httpSwagger.PersistAuthorization(true),
		))
	}

	// Profiling endpoints (if enabled)
	if r.config.EnableProfiling {
		r.mux.Mount("/debug", middleware.Profiler())
	}

	// API v1 routes
	r.mux.Route("/v1", func(v1 chi.Router) {
		// Authentication routes (some unauthenticated)
		v1.Route("/auth", func(auth chi.Router) {
			auth.Post("/sync", r.authHandler.SyncUser)
		})

		// Protected routes (require authentication)
		v1.Group(func(protected chi.Router) {
			protected.Use(r.AuthMiddleware)

			// User profile
			protected.Get("/me", r.userHandler.GetMe)
			protected.Patch("/me", r.userHandler.UpdateMe)

			// Experiences
			protected.Route("/experiences", func(exp chi.Router) {
				exp.Get("/", r.experienceHandler.List)
				exp.Post("/", r.experienceHandler.Create)

				exp.Route("/{experienceID}", func(expByID chi.Router) {
					expByID.Get("/", r.experienceHandler.Get)
					expByID.Put("/", r.experienceHandler.Update)
					expByID.Delete("/", r.experienceHandler.Delete)

					// Bullets under experience
					expByID.Post("/bullets", r.bulletHandler.Create)
				})
			})

			// Bullets (direct access)
			protected.Route("/bullets", func(bullet chi.Router) {
				bullet.Route("/{bulletID}", func(bulletByID chi.Router) {
					bulletByID.Put("/", r.bulletHandler.Update)
					bulletByID.Delete("/", r.bulletHandler.Delete)
					bulletByID.Post("/score", r.bulletHandler.RecalculateScore)
				})
			})

			// Skills
			protected.Route("/skills", func(skill chi.Router) {
				skill.Get("/", r.skillHandler.List)
				skill.Post("/batch", r.skillHandler.BatchUpsert)

				skill.Route("/{skillID}", func(skillByID chi.Router) {
					skillByID.Delete("/", r.skillHandler.Delete)
				})
			})

			// Spoken Languages
			protected.Route("/languages", func(lang chi.Router) {
				lang.Get("/", r.languageHandler.List)
				lang.Post("/", r.languageHandler.Create)

				lang.Route("/{languageID}", func(langByID chi.Router) {
					langByID.Delete("/", r.languageHandler.Delete)
				})
			})

			// Resumes
			protected.Route("/resumes", func(resume chi.Router) {
				resume.Get("/", r.resumeHandler.List)
				resume.Post("/", r.resumeHandler.Create)

				resume.Route("/{resumeID}", func(resumeByID chi.Router) {
					resumeByID.Get("/", r.resumeHandler.Get)
					resumeByID.Delete("/", r.resumeHandler.Delete)
					resumeByID.Post("/tailor", r.resumeHandler.Tailor)
					resumeByID.Patch("/content", r.resumeHandler.UpdateStatus)
					resumeByID.Get("/pdf", r.resumeHandler.GeneratePDF)
				})
			})

			// Tools
			protected.Route("/tools", func(tools chi.Router) {
				tools.Post("/parse-job", r.toolsHandler.ParseJobURL)
			})
		})
	})

	// Not found handler
	r.mux.NotFound(r.notFoundHandler)

	// Method not allowed handler
	r.mux.MethodNotAllowed(r.methodNotAllowedHandler)
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// healthHandler returns the health status of the service.
//
//	@Summary		Health check
//	@Description	Returns the health status of the service
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	HealthResponse
//	@Router			/health [get]
func (r *Router) healthHandler(w http.ResponseWriter, req *http.Request) {
	respondJSON(w, http.StatusOK, HealthResponse{
		Status:  "healthy",
		Service: "chameleon-vitae",
	})
}

// notFoundHandler handles 404 errors.
func (r *Router) notFoundHandler(w http.ResponseWriter, req *http.Request) {
	respondError(w, http.StatusNotFound, "RESOURCE_NOT_FOUND", "The requested resource was not found")
}

// methodNotAllowedHandler handles 405 errors.
func (r *Router) methodNotAllowedHandler(w http.ResponseWriter, req *http.Request) {
	respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not allowed for this resource")
}
