// Package config provides application configuration using Viper.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Firebase FirebaseConfig
	Groq     GroqConfig
	Jina     JinaConfig
	PDF      PDFConfig
	Storage  StorageConfig
}

// AppConfig contains general application settings.
type AppConfig struct {
	Name        string
	Version     string
	Environment string
	Debug       bool
}

// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	Port            int
	Host            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	MaxRequestSize  int64
	AllowedOrigins  []string
	EnableSwagger   bool
	EnableProfiling bool
}

// DatabaseConfig contains PostgreSQL connection settings.
type DatabaseConfig struct {
	Host              string
	Port              int
	User              string
	Password          string
	Database          string
	SSLMode           string
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetime   time.Duration
	ConnMaxIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

// FirebaseConfig contains Firebase authentication settings.
type FirebaseConfig struct {
	ProjectID       string
	CredentialsFile string
	CredentialsJSON string
}

// GroqConfig contains Groq AI provider settings.
type GroqConfig struct {
	APIKey         string
	BaseURL        string
	DefaultModel   string
	AnalysisModel  string
	MaxRetries     int
	RequestTimeout time.Duration
}

// JinaConfig contains Jina Reader settings.
type JinaConfig struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
}

// PDFConfig contains Gotenberg PDF engine settings.
type PDFConfig struct {
	BaseURL string
	Timeout time.Duration
}

// StorageConfig contains file storage settings.
type StorageConfig struct {
	Type      string // "local" or "s3" or "gcs"
	LocalPath string
	S3Bucket  string
	S3Region  string
}

// Load loads configuration from environment variables and config files.
func Load() (*Config, error) {
	v := viper.New()

	// Set config file options
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/chameleon-vitae")

	// Enable environment variable reading
	v.SetEnvPrefix("CHAMELEON")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Set defaults
	setDefaults(v)

	// Read config file (optional)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found is acceptable - we use env vars
	}

	cfg := &Config{}
	if err := unmarshalConfig(v, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// setDefaults sets default configuration values.
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "chameleon-vitae")
	v.SetDefault("app.version", "0.1.0")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.debug", false)

	// Server defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.readTimeout", "15s")
	v.SetDefault("server.writeTimeout", "30s")
	v.SetDefault("server.idleTimeout", "60s")
	v.SetDefault("server.maxRequestSize", 10*1024*1024) // 10MB
	v.SetDefault("server.allowedOrigins", []string{"*"})
	v.SetDefault("server.enableSwagger", true)
	v.SetDefault("server.enableProfiling", false)

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "chameleon")
	v.SetDefault("database.password", "")
	v.SetDefault("database.database", "chameleon_vitae")
	v.SetDefault("database.sslMode", "disable")
	v.SetDefault("database.maxOpenConns", 25)
	v.SetDefault("database.maxIdleConns", 5)
	v.SetDefault("database.connMaxLifetime", "30m")
	v.SetDefault("database.connMaxIdleTime", "5m")
	v.SetDefault("database.healthCheckPeriod", "1m")

	// Firebase defaults
	v.SetDefault("firebase.projectId", "")
	v.SetDefault("firebase.credentialsFile", "")
	v.SetDefault("firebase.credentialsJson", "")

	// Groq defaults
	v.SetDefault("groq.apiKey", "")
	v.SetDefault("groq.baseUrl", "https://api.groq.com/openai/v1")
	v.SetDefault("groq.defaultModel", "llama-3.3-70b-versatile")
	v.SetDefault("groq.analysisModel", "llama-4-scout-17b-16e-instruct")
	v.SetDefault("groq.maxRetries", 3)
	v.SetDefault("groq.requestTimeout", "60s")

	// Jina defaults
	v.SetDefault("jina.apiKey", "")
	v.SetDefault("jina.baseUrl", "https://r.jina.ai")
	v.SetDefault("jina.timeout", "30s")

	// PDF defaults
	v.SetDefault("pdf.baseUrl", "http://localhost:3000")
	v.SetDefault("pdf.timeout", "60s")

	// Storage defaults
	v.SetDefault("storage.type", "local")
	v.SetDefault("storage.localPath", "./storage")
	v.SetDefault("storage.s3Bucket", "")
	v.SetDefault("storage.s3Region", "")
}

// unmarshalConfig unmarshals viper config into the Config struct.
func unmarshalConfig(v *viper.Viper, cfg *Config) error {
	// App
	cfg.App.Name = v.GetString("app.name")
	cfg.App.Version = v.GetString("app.version")
	cfg.App.Environment = v.GetString("app.environment")
	cfg.App.Debug = v.GetBool("app.debug")

	// Server
	cfg.Server.Port = v.GetInt("server.port")
	cfg.Server.Host = v.GetString("server.host")
	cfg.Server.ReadTimeout = v.GetDuration("server.readTimeout")
	cfg.Server.WriteTimeout = v.GetDuration("server.writeTimeout")
	cfg.Server.IdleTimeout = v.GetDuration("server.idleTimeout")
	cfg.Server.MaxRequestSize = v.GetInt64("server.maxRequestSize")
	cfg.Server.AllowedOrigins = v.GetStringSlice("server.allowedOrigins")
	cfg.Server.EnableSwagger = v.GetBool("server.enableSwagger")
	cfg.Server.EnableProfiling = v.GetBool("server.enableProfiling")

	// Database
	cfg.Database.Host = v.GetString("database.host")
	cfg.Database.Port = v.GetInt("database.port")
	cfg.Database.User = v.GetString("database.user")
	cfg.Database.Password = v.GetString("database.password") // pragma: allowlist secret
	cfg.Database.Database = v.GetString("database.database")
	cfg.Database.SSLMode = v.GetString("database.sslMode")
	cfg.Database.MaxOpenConns = v.GetInt("database.maxOpenConns")
	cfg.Database.MaxIdleConns = v.GetInt("database.maxIdleConns")
	cfg.Database.ConnMaxLifetime = v.GetDuration("database.connMaxLifetime")
	cfg.Database.ConnMaxIdleTime = v.GetDuration("database.connMaxIdleTime")
	cfg.Database.HealthCheckPeriod = v.GetDuration("database.healthCheckPeriod")

	// Firebase
	cfg.Firebase.ProjectID = v.GetString("firebase.projectId")
	cfg.Firebase.CredentialsFile = v.GetString("firebase.credentialsFile")
	cfg.Firebase.CredentialsJSON = v.GetString("firebase.credentialsJson")

	// Groq
	cfg.Groq.APIKey = v.GetString("groq.apiKey") // pragma: allowlist secret
	cfg.Groq.BaseURL = v.GetString("groq.baseUrl")
	cfg.Groq.DefaultModel = v.GetString("groq.defaultModel")
	cfg.Groq.AnalysisModel = v.GetString("groq.analysisModel")
	cfg.Groq.MaxRetries = v.GetInt("groq.maxRetries")
	cfg.Groq.RequestTimeout = v.GetDuration("groq.requestTimeout")

	// Jina
	cfg.Jina.APIKey = v.GetString("jina.apiKey") // pragma: allowlist secret
	cfg.Jina.BaseURL = v.GetString("jina.baseUrl")
	cfg.Jina.Timeout = v.GetDuration("jina.timeout")

	// PDF
	cfg.PDF.BaseURL = v.GetString("pdf.baseUrl")
	cfg.PDF.Timeout = v.GetDuration("pdf.timeout")

	// Storage
	cfg.Storage.Type = v.GetString("storage.type")
	cfg.Storage.LocalPath = v.GetString("storage.localPath")
	cfg.Storage.S3Bucket = v.GetString("storage.s3Bucket")
	cfg.Storage.S3Region = v.GetString("storage.s3Region")

	return nil
}

// validateConfig validates required configuration fields.
func validateConfig(cfg *Config) error {
	// Firebase is required for authentication
	if cfg.Firebase.ProjectID == "" {
		return fmt.Errorf("firebase.projectId is required")
	}

	// Groq API key is required for AI features
	if cfg.Groq.APIKey == "" {
		return fmt.Errorf("groq.apiKey is required")
	}

	// Database password should be set in production
	if cfg.App.Environment == "production" && cfg.Database.Password == "" {
		return fmt.Errorf("database.password is required in production")
	}

	return nil
}

// DSN returns the PostgreSQL connection string.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s", // pragma: allowlist secret
		c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode,
	)
}

// IsProduction returns true if running in production environment.
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// IsDevelopment returns true if running in development environment.
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}
