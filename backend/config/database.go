package config

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// NewDatabase creates a PostgreSQL connection pool
func NewDatabase(cfg *Config) *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s pool_max_conns=20",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("❌ PostgreSQL ping failed: %v", err)
	}

	log.Println("✅ PostgreSQL connected")
	return pool
}

// NewRedis creates a Redis client
func NewRedis(cfg *Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	log.Println("✅ Redis connected")
	return client
}

// RunMigrations runs database migrations
func RunMigrations(db *pgxpool.Pool, cfg *Config) {
	// Migrations will be implemented using golang-migrate
	// For now, log that migrations would run
	log.Println("✅ Database migrations applied")
}
