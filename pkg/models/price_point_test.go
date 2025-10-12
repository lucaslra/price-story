package models

import (
    "encoding/json"
    "testing"
    "time"
)

func TestPricePointJSON_OmitsUpdatedByUserWhenNil(t *testing.T) {
    pp := PricePoint{
        ID:        "pp-1",
        Price:     1.23,
        Timestamp: time.Now(),
        Product:   Product{ID: "prod-1", ProductName: "Widget"},
        CreatedByUser: User{ID: "u1", Email: "u1@example.com"},
        UpdatedByUser: nil,
    }

    b, err := json.Marshal(pp)
    if err != nil {
        t.Fatalf("marshal failed: %v", err)
    }

    var m map[string]any
    if err := json.Unmarshal(b, &m); err != nil {
        t.Fatalf("unmarshal failed: %v", err)
    }

    if _, ok := m["updated_by_user"]; ok {
        t.Fatalf("expected updated_by_user to be omitted when nil")
    }
    if m["price"] != 1.23 {
        t.Fatalf("unexpected price: %v", m["price"])
    }
    if _, ok := m["product"].(map[string]any); !ok {
        t.Fatalf("expected product to be an object")
    }
}