package models

import (
	"encoding/json"
	"testing"
)

func TestProductJSON_OmitsUpdatedByUserWhenNil(t *testing.T) {
	p := Product{
		ID:                 "prod-1",
		ProductName:        "Widget",
		ProductUrl:         "https://example.com/widget",
		ProductImageUrl:    "https://example.com/widget.png",
		ProductDescription: "A useful widget",
		CreatedByUser:      User{ID: "u1", Email: "u1@example.com"},
		UpdatedByUser:      nil,
	}

	b, err := json.Marshal(p)
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

	if got := m["product_url"]; got != "https://example.com/widget" {
		t.Fatalf("unexpected product_url: %v", got)
	}
	if got := m["product_image_url"]; got != "https://example.com/widget.png" {
		t.Fatalf("unexpected product_image_url: %v", got)
	}
}

func TestProductJSON_IncludesUpdatedByUserWhenPresent(t *testing.T) {
	updater := &User{ID: "u2", Email: "u2@example.com"}
	p := Product{
		ID:                 "prod-2",
		ProductName:        "Gadget",
		ProductUrl:         "https://example.com/gadget",
		ProductImageUrl:    "https://example.com/gadget.png",
		ProductDescription: "A fancy gadget",
		CreatedByUser:      User{ID: "u1", Email: "u1@example.com"},
		UpdatedByUser:      updater,
	}

	b, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	v, ok := m["updated_by_user"]
	if !ok {
		t.Fatalf("expected updated_by_user to be present when non-nil")
	}
	nested, ok := v.(map[string]any)
	if !ok {
		t.Fatalf("expected updated_by_user to be an object, got %T", v)
	}
	if nested["id"] != "u2" {
		t.Fatalf("unexpected updated_by_user.id: %v", nested["id"])
	}
}
