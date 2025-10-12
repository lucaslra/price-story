package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

func TestHealthCheckHandler_OK(t *testing.T) {
    h := &Handler{}
    rr := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/api/health", nil)
    h.HealthCheckHandler(rr, req)
    if rr.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", rr.Code)
    }
}

func TestCreateProductHandler_InvalidURL(t *testing.T) {
    h := &Handler{}
    payload := map[string]any{
        "product_name": "Widget",
        "product_url":  "ftp://invalid",
    }
    b, _ := json.Marshal(payload)
    rr := httptest.NewRecorder()
    req := httptest.NewRequest("POST", "/api/products", bytes.NewReader(b))
    h.CreateProductHandler(rr, req)
    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected 400, got %d", rr.Code)
    }
}

func TestCreatePricePointHandler_InvalidPrice(t *testing.T) {
    h := &Handler{}
    payload := map[string]any{
        "price":      0.0,
        "timestamp":  time.Now(),
        "product_id": "123e4567-e89b-12d3-a456-426614174000",
    }
    b, _ := json.Marshal(payload)
    rr := httptest.NewRecorder()
    req := httptest.NewRequest("POST", "/api/price-points", bytes.NewReader(b))
    h.CreatePricePointHandler(rr, req)
    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected 400, got %d", rr.Code)
    }
}