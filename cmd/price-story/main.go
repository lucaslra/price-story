package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"price-story/pkg/config"
	"price-story/pkg/database"
	"price-story/pkg/migrations"
	"price-story/pkg/router"
)

//go:embed migrations
var migrationsFS embed.FS

//go:embed openapi/*
//go:embed static/swagger/*
var docsFS embed.FS

func main() {
	// Parse command line flags
	var migrateCmd = flag.String("migrate", "", "Migration command: up, down, status")
	flag.Parse()

	// Load configuration
	cfg := config.New()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Successfully connected to database")

	// Handle migration commands
	if *migrateCmd != "" {
		migrator := migrations.NewMigrator(db, migrationsFS)
		if err := migrator.LoadMigrations(); err != nil {
			log.Fatalf("Failed to load migrations: %v", err)
		}

		switch *migrateCmd {
		case "up":
			if err := migrator.Up(); err != nil {
				log.Fatalf("Failed to run migrations: %v", err)
			}
		case "down":
			if err := migrator.Down(); err != nil {
				log.Fatalf("Failed to rollback migration: %v", err)
			}
		case "status":
			if err := migrator.Status(); err != nil {
				log.Fatalf("Failed to get migration status: %v", err)
			}
		default:
			fmt.Printf("Unknown migration command: %s\n", *migrateCmd)
			fmt.Println("Available commands: up, down, status")
			os.Exit(1)
		}
		return
	}

	// Run migrations automatically
	migrator := migrations.NewMigrator(db, migrationsFS)
	if err := migrator.LoadMigrations(); err != nil {
		log.Fatalf("Failed to load migrations: %v", err)
	}
	if err := migrator.Up(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

	// Initialize router with database connection
	r := router.New(db)

	// Serve OpenAPI spec
	r.HandleFunc("/api/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		b, err := docsFS.ReadFile("openapi/openapi.yaml")
		if err != nil {
			http.Error(w, "OpenAPI spec not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/yaml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	})

	// Serve Swagger UI
	r.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		b, err := docsFS.ReadFile("static/swagger/index.html")
		if err != nil {
			http.Error(w, "Swagger UI not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	})

	// Configure server
	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + cfg.Server.Port,
		WriteTimeout: cfg.Server.WriteTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
	}

	// Start server
	fmt.Printf("Server starting on port %s...\n", cfg.Server.Port)
	log.Fatal(srv.ListenAndServe())
}
