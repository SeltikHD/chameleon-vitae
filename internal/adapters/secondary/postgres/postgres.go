// Package postgres provides PostgreSQL database adapters implementing repository interfaces.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config holds PostgreSQL connection configuration.
type Config struct {
	// Host is the database host.
	Host string

	// Port is the database port.
	Port int

	// User is the database user.
	User string

	// Password is the database password.
	Password string

	// Database is the database name.
	Database string

	// SSLMode is the SSL mode (disable, require, verify-ca, verify-full).
	SSLMode string

	// MaxConns is the maximum number of connections in the pool.
	MaxConns int

	// MinConns is the minimum number of connections in the pool.
	MinConns int

	// MaxConnLifetime is the maximum lifetime of a connection.
	MaxConnLifetime time.Duration

	// MaxConnIdleTime is the maximum idle time for a connection.
	MaxConnIdleTime time.Duration

	// HealthCheckPeriod is the period between health checks.
	HealthCheckPeriod time.Duration

	// ConnectionURL is an optional full connection URL (overrides other fields).
	ConnectionURL string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Host:              "localhost",
		Port:              5432,
		User:              "chameleon",
		Database:          "chameleon_vitae",
		SSLMode:           "disable",
		MaxConns:          25,
		MinConns:          5,
		MaxConnLifetime:   time.Hour,
		MaxConnIdleTime:   30 * time.Minute,
		HealthCheckPeriod: time.Minute,
	}
}

// DSN returns the connection string for the database.
func (c Config) DSN() string {
	if c.ConnectionURL != "" {
		return c.ConnectionURL
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
	)
}

// DB wraps a pgxpool.Pool and provides repository factories.
type DB struct {
	pool *pgxpool.Pool
}

// New creates a new DB connection pool.
func New(ctx context.Context, cfg Config) (*DB, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Apply pool settings.
	poolCfg.MaxConns = int32(cfg.MaxConns)
	poolCfg.MinConns = int32(cfg.MinConns)
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = cfg.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection.
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{pool: pool}, nil
}

// Close closes the database connection pool.
func (db *DB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

// Pool returns the underlying connection pool for advanced usage.
func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

// HealthCheck verifies the database connection is alive.
func (db *DB) HealthCheck(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// UserRepository returns a new UserRepository instance.
func (db *DB) UserRepository() *UserRepository {
	return &UserRepository{pool: db.pool}
}

// ExperienceRepository returns a new ExperienceRepository instance.
func (db *DB) ExperienceRepository() *ExperienceRepository {
	return &ExperienceRepository{pool: db.pool}
}

// BulletRepository returns a new BulletRepository instance.
func (db *DB) BulletRepository() *BulletRepository {
	return &BulletRepository{pool: db.pool}
}

// SkillRepository returns a new SkillRepository instance.
func (db *DB) SkillRepository() *SkillRepository {
	return &SkillRepository{pool: db.pool}
}

// SpokenLanguageRepository returns a new SpokenLanguageRepository instance.
func (db *DB) SpokenLanguageRepository() *SpokenLanguageRepository {
	return &SpokenLanguageRepository{pool: db.pool}
}

// ResumeRepository returns a new ResumeRepository instance.
func (db *DB) ResumeRepository() *ResumeRepository {
	return &ResumeRepository{pool: db.pool}
}

// EducationRepository returns a new EducationRepository instance.
func (db *DB) EducationRepository() *EducationRepository {
	return &EducationRepository{pool: db.pool}
}

// ProjectRepository returns a new ProjectRepository instance.
func (db *DB) ProjectRepository() *ProjectRepository {
	return &ProjectRepository{pool: db.pool}
}

// ProjectBulletRepository returns a new ProjectBulletRepository instance.
func (db *DB) ProjectBulletRepository() *ProjectBulletRepository {
	return &ProjectBulletRepository{pool: db.pool}
}
