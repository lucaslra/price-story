package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"price-story/pkg/middleware"
	"price-story/pkg/models"
	"price-story/pkg/repository"
	"price-story/pkg/validation"
	"time"

	"github.com/gorilla/mux"
)

// UUID validation centralized in pkg/validation

// Handler struct holds dependencies for handlers
type Handler struct {
	UserRepo       *repository.UserRepository
	ProductRepo    *repository.ProductRepository
	PriceStoryRepo *repository.PriceStoryRepository
	PricePointRepo *repository.PricePointRepository
}

// NewHandler creates a new handler with dependencies
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		UserRepo:       repository.NewUserRepository(db),
		ProductRepo:    repository.NewProductRepository(db),
		PriceStoryRepo: repository.NewPriceStoryRepository(db),
		PricePointRepo: repository.NewPricePointRepository(db),
	}
}

// HealthCheckHandler returns the API health status
func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// ----- Users -----

// GetUsersHandler returns a list of users
func (h *Handler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := h.UserRepo.ListUsers()
	if err != nil {
		log.Printf("Error getting users: %v", err)
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"users": users}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// GetUserHandler returns a specific user by ID
func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	if !validation.IsUUID(id) {
		http.Error(w, "id must be a valid UUID", http.StatusBadRequest)
		return
	}
	user, err := h.UserRepo.GetUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			log.Printf("Error getting user %s: %v", id, err)
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// CreateUserHandler creates a new user
func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if u.Email == "" || u.PasswordHash == "" {
		http.Error(w, "Email and passwordHash are required", http.StatusBadRequest)
		return
	}
	created, err := h.UserRepo.CreateUser(u)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// UpdateUserHandler updates an existing user
func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	if !validation.IsUUID(id) {
		http.Error(w, "id must be a valid UUID", http.StatusBadRequest)
		return
	}
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if u.Email == "" || u.PasswordHash == "" {
		http.Error(w, "Email and passwordHash are required", http.StatusBadRequest)
		return
	}
	updated, err := h.UserRepo.UpdateUser(id, u)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			log.Printf("Error updating user %s: %v", id, err)
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// DeleteUserHandler deletes a user
func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if !validation.IsUUID(id) {
		http.Error(w, "id must be a valid UUID", http.StatusBadRequest)
		return
	}
	deleted, err := h.UserRepo.DeleteUser(id)
	if err != nil {
		log.Printf("Error deleting user %s: %v", id, err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ----- Products -----

// GetProductsHandler returns a list of products
func (h *Handler) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	products, err := h.ProductRepo.ListProducts()
	if err != nil {
		log.Printf("Error getting products: %v", err)
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"products": products}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// GetProductHandler returns a specific product by ID
func (h *Handler) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	if !validation.IsUUID(id) {
		http.Error(w, "id must be a valid UUID", http.StatusBadRequest)
		return
	}
	product, err := h.ProductRepo.GetProduct(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			log.Printf("Error getting product %s: %v", id, err)
			http.Error(w, "Failed to retrieve product", http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(product); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// CreateProductHandler creates a new product
func (h *Handler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Accept product fields only; created/updated users are filled internally
	var payload struct {
		ProductName        string `json:"product_name"`
		ProductUrl         string `json:"product_url"`
		ProductImageUrl    string `json:"product_image_url"`
		ProductDescription string `json:"product_description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if payload.ProductName == "" || payload.ProductUrl == "" {
		http.Error(w, "productName and productUrl are required", http.StatusBadRequest)
		return
	}
	// URL validation for product and optional image URL
	if !validation.IsValidURL(payload.ProductUrl) {
		http.Error(w, "productUrl must be a valid http(s) URL", http.StatusBadRequest)
		return
	}
	if payload.ProductImageUrl != "" && !validation.IsValidURL(payload.ProductImageUrl) {
		http.Error(w, "productImageUrl must be a valid http(s) URL", http.StatusBadRequest)
		return
	}
	// Fill created/updated users from middleware-injected actor
	actorID, ok := middleware.ActorUserIDFromContext(r.Context())
	if !ok || actorID == "" {
		http.Error(w, "Actor user not available", http.StatusInternalServerError)
		return
	}
	created, err := h.ProductRepo.CreateProduct(models.Product{
		ProductName:        payload.ProductName,
		ProductUrl:         payload.ProductUrl,
		ProductImageUrl:    payload.ProductImageUrl,
		ProductDescription: payload.ProductDescription,
		CreatedByUser:      models.User{ID: actorID},
		UpdatedByUser:      nil,
	})
	if err != nil {
		log.Printf("Error creating product: %v", err)
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// UpdateProductHandler updates an existing product
func (h *Handler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	if !validation.IsUUID(id) {
		http.Error(w, "id must be a valid UUID", http.StatusBadRequest)
		return
	}
	// Accept product fields only; updated user is filled internally
	var payload struct {
		ProductName        string `json:"product_name"`
		ProductUrl         string `json:"product_url"`
		ProductImageUrl    string `json:"product_image_url"`
		ProductDescription string `json:"product_description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if payload.ProductName == "" || payload.ProductUrl == "" {
		http.Error(w, "productName and productUrl are required", http.StatusBadRequest)
		return
	}
	// URL validation for product and optional image URL
	if !validation.IsValidURL(payload.ProductUrl) {
		http.Error(w, "productUrl must be a valid http(s) URL", http.StatusBadRequest)
		return
	}
	if payload.ProductImageUrl != "" && !validation.IsValidURL(payload.ProductImageUrl) {
		http.Error(w, "productImageUrl must be a valid http(s) URL", http.StatusBadRequest)
		return
	}
	actorID, ok := middleware.ActorUserIDFromContext(r.Context())
	if !ok || actorID == "" {
		http.Error(w, "Actor user not available", http.StatusInternalServerError)
		return
	}
	updated, err := h.ProductRepo.UpdateProduct(id, models.Product{
		ProductName:        payload.ProductName,
		ProductUrl:         payload.ProductUrl,
		ProductImageUrl:    payload.ProductImageUrl,
		ProductDescription: payload.ProductDescription,
		UpdatedByUser:      &models.User{ID: actorID},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			log.Printf("Error updating product %s: %v", id, err)
			http.Error(w, "Failed to update product", http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// DeleteProductHandler deletes a product
func (h *Handler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if !validation.IsUUID(id) {
		http.Error(w, "id must be a valid UUID", http.StatusBadRequest)
		return
	}
	deleted, err := h.ProductRepo.DeleteProduct(id)
	if err != nil {
		log.Printf("Error deleting product %s: %v", id, err)
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ----- Price Stories -----

// GetPriceStoriesHandler returns a list of price stories
func (h *Handler) GetPriceStoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stories, err := h.PriceStoryRepo.ListPriceStories()
	if err != nil {
		log.Printf("Error getting price stories: %v", err)
		http.Error(w, "Failed to retrieve price stories", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"priceStories": stories}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// GetPriceStoryHandler returns a specific price story by ID
func (h *Handler) GetPriceStoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	if !validation.IsUUID(id) {
		http.Error(w, "id must be a valid UUID", http.StatusBadRequest)
		return
	}
	ps, err := h.PriceStoryRepo.GetPriceStory(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Price story not found", http.StatusNotFound)
		} else {
			log.Printf("Error getting price story %s: %v", id, err)
			http.Error(w, "Failed to retrieve price story", http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(ps); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// CreatePriceStoryHandler creates a new price story
func (h *Handler) CreatePriceStoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Accept product_id only; created/updated users are filled internally
	var payload struct {
		ProductID string `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if payload.ProductID == "" {
		http.Error(w, "productId is required", http.StatusBadRequest)
		return
	}
	if !validation.IsUUID(payload.ProductID) {
		http.Error(w, "productId must be a valid UUID", http.StatusBadRequest)
		return
	}
	actorID, ok := middleware.ActorUserIDFromContext(r.Context())
	if !ok || actorID == "" {
		http.Error(w, "Actor user not available", http.StatusInternalServerError)
		return
	}
	created, err := h.PriceStoryRepo.CreatePriceStory(models.PriceStory{
		Product:       models.Product{ID: payload.ProductID},
		CreatedByUser: models.User{ID: actorID},
		UpdatedByUser: nil,
	})
	if err != nil {
		log.Printf("Error creating price story: %v", err)
		http.Error(w, "Failed to create price story", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// UpdatePriceStoryHandler updates an existing price story
func (h *Handler) UpdatePriceStoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	if !validation.IsUUID(id) {
		http.Error(w, "id must be a valid UUID", http.StatusBadRequest)
		return
	}
	// Accept product_id only; updated user is filled internally
	var payload struct {
		ProductID string `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if payload.ProductID == "" {
		http.Error(w, "productId is required", http.StatusBadRequest)
		return
	}
	if !validation.IsUUID(payload.ProductID) {
		http.Error(w, "productId must be a valid UUID", http.StatusBadRequest)
		return
	}
	actorID, ok := middleware.ActorUserIDFromContext(r.Context())
	if !ok || actorID == "" {
		http.Error(w, "Actor user not available", http.StatusInternalServerError)
		return
	}
	updated, err := h.PriceStoryRepo.UpdatePriceStory(id, models.PriceStory{
		Product:       models.Product{ID: payload.ProductID},
		UpdatedByUser: &models.User{ID: actorID},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Price story not found", http.StatusNotFound)
		} else {
			log.Printf("Error updating price story %s: %v", id, err)
			http.Error(w, "Failed to update price story", http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// DeletePriceStoryHandler deletes a price story
func (h *Handler) DeletePriceStoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if !validation.IsUUID(id) {
		http.Error(w, "id must be a valid UUID", http.StatusBadRequest)
		return
	}
	deleted, err := h.PriceStoryRepo.DeletePriceStory(id)
	if err != nil {
		log.Printf("Error deleting price story %s: %v", id, err)
		http.Error(w, "Failed to delete price story", http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(w, "Price story not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ----- Price Points -----

// GetPricePointsHandler returns a list of price points
func (h *Handler) GetPricePointsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	points, err := h.PricePointRepo.ListPricePoints()
	if err != nil {
		log.Printf("Error getting price points: %v", err)
		http.Error(w, "Failed to retrieve price points", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"pricePoints": points}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// GetPricePointHandler returns a specific price point by ID
func (h *Handler) GetPricePointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	if !validation.IsUUID(id) {
		http.Error(w, "id must be a valid UUID", http.StatusBadRequest)
		return
	}
	pp, err := h.PricePointRepo.GetPricePoint(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Price point not found", http.StatusNotFound)
		} else {
			log.Printf("Error getting price point %s: %v", id, err)
			http.Error(w, "Failed to retrieve price point", http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(pp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// CreatePricePointHandler creates a new price point
func (h *Handler) CreatePricePointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Accept product_id, price, timestamp only; created/updated users are filled internally
	var payload struct {
		Price     float64   `json:"price"`
		Timestamp time.Time `json:"timestamp"`
		ProductID string    `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// timestamp comes as time.Time in current model; ensure presence
	if !validation.IsValidPrice(payload.Price) || payload.ProductID == "" || payload.Timestamp.IsZero() {
		http.Error(w, "price, timestamp, productId are required", http.StatusBadRequest)
		return
	}
	if !validation.IsUUID(payload.ProductID) {
		http.Error(w, "productId must be a valid UUID", http.StatusBadRequest)
		return
	}
	actorID, ok := middleware.ActorUserIDFromContext(r.Context())
	if !ok || actorID == "" {
		http.Error(w, "Actor user not available", http.StatusInternalServerError)
		return
	}
	created, err := h.PricePointRepo.CreatePricePoint(models.PricePoint{
		Price:         payload.Price,
		Timestamp:     payload.Timestamp,
		Product:       models.Product{ID: payload.ProductID},
		CreatedByUser: models.User{ID: actorID},
		UpdatedByUser: nil,
	})
	if err != nil {
		log.Printf("Error creating price point: %v", err)
		http.Error(w, "Failed to create price point", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// UpdatePricePointHandler updates an existing price point
func (h *Handler) UpdatePricePointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	if !validation.IsUUID(id) {
		http.Error(w, "id must be a valid UUID", http.StatusBadRequest)
		return
	}
	// Accept product_id, price, timestamp only; updated user is filled internally
	var payload struct {
		Price     float64   `json:"price"`
		Timestamp time.Time `json:"timestamp"`
		ProductID string    `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if !validation.IsValidPrice(payload.Price) || payload.ProductID == "" || payload.Timestamp.IsZero() {
		http.Error(w, "price, timestamp, productId are required", http.StatusBadRequest)
		return
	}
	if !validation.IsUUID(payload.ProductID) {
		http.Error(w, "productId must be a valid UUID", http.StatusBadRequest)
		return
	}
	actorID, ok := middleware.ActorUserIDFromContext(r.Context())
	if !ok || actorID == "" {
		http.Error(w, "Actor user not available", http.StatusInternalServerError)
		return
	}
	updated, err := h.PricePointRepo.UpdatePricePoint(id, models.PricePoint{
		Price:         payload.Price,
		Timestamp:     payload.Timestamp,
		Product:       models.Product{ID: payload.ProductID},
		UpdatedByUser: &models.User{ID: actorID},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Price point not found", http.StatusNotFound)
		} else {
			log.Printf("Error updating price point %s: %v", id, err)
			http.Error(w, "Failed to update price point", http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// DeletePricePointHandler deletes a price point
func (h *Handler) DeletePricePointHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if !validation.IsUUID(id) {
		http.Error(w, "id must be a valid UUID", http.StatusBadRequest)
		return
	}
	deleted, err := h.PricePointRepo.DeletePricePoint(id)
	if err != nil {
		log.Printf("Error deleting price point %s: %v", id, err)
		http.Error(w, "Failed to delete price point", http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(w, "Price point not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
