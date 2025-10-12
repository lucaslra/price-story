package router

import (
	"database/sql"
	"log"
	"price-story/pkg/handlers"
	"price-story/pkg/middleware"
	"price-story/pkg/repository"

	"github.com/gorilla/mux"
)

// New creates and configures a new router with all application routes
func New(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// Apply middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Ensure a system placeholder user exists and inject its ID into context
	ur := repository.NewUserRepository(db)
	sys, err := ur.EnsureUserByEmail("system@example.com", "placeholder")
	if err != nil {
		log.Printf("Failed to ensure system user: %v", err)
	}
	r.Use(middleware.InjectActorUser(sys.ID))

	// Initialize handlers with database connection
	h := handlers.NewHandler(db)

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Health check endpoint
	api.HandleFunc("/health", h.HealthCheckHandler).Methods("GET")

	// User endpoints
	api.HandleFunc("/users", h.GetUsersHandler).Methods("GET")
	api.HandleFunc("/users/{id}", h.GetUserHandler).Methods("GET")
	api.HandleFunc("/users", h.CreateUserHandler).Methods("POST")
	api.HandleFunc("/users/{id}", h.UpdateUserHandler).Methods("PUT")
	api.HandleFunc("/users/{id}", h.DeleteUserHandler).Methods("DELETE")

	// Product endpoints
	api.HandleFunc("/products", h.GetProductsHandler).Methods("GET")
	api.HandleFunc("/products/{id}", h.GetProductHandler).Methods("GET")
	api.HandleFunc("/products", h.CreateProductHandler).Methods("POST")
	api.HandleFunc("/products/{id}", h.UpdateProductHandler).Methods("PUT")
	api.HandleFunc("/products/{id}", h.DeleteProductHandler).Methods("DELETE")

	// Price Story endpoints
	api.HandleFunc("/price-stories", h.GetPriceStoriesHandler).Methods("GET")
	api.HandleFunc("/price-stories/{id}", h.GetPriceStoryHandler).Methods("GET")
	api.HandleFunc("/price-stories", h.CreatePriceStoryHandler).Methods("POST")
	api.HandleFunc("/price-stories/{id}", h.UpdatePriceStoryHandler).Methods("PUT")
	api.HandleFunc("/price-stories/{id}", h.DeletePriceStoryHandler).Methods("DELETE")

	// Price Point endpoints
	api.HandleFunc("/price-points", h.GetPricePointsHandler).Methods("GET")
	api.HandleFunc("/price-points/{id}", h.GetPricePointHandler).Methods("GET")
	api.HandleFunc("/price-points", h.CreatePricePointHandler).Methods("POST")
	api.HandleFunc("/price-points/{id}", h.UpdatePricePointHandler).Methods("PUT")
	api.HandleFunc("/price-points/{id}", h.DeletePricePointHandler).Methods("DELETE")

	return r
}
