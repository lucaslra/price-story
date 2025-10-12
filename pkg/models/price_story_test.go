package models

import (
	"encoding/json"
	"testing"
)

func TestPriceStoryJSON_OmitsUpdatedByUserWhenNil(t *testing.T) {
	ps := PriceStory{
		ID:            "ps-1",
		Product:       Product{ID: "prod-1", ProductName: "Widget"},
		CreatedByUser: User{ID: "u1", Email: "u1@example.com"},
		UpdatedByUser: nil,
	}

	b, err := json.Marshal(ps)
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
	if _, ok := m["product"].(map[string]any); !ok {
		t.Fatalf("expected product to be an object")
	}
}
