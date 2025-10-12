package openapi

import (
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"

	"price-story/pkg/queries"
	"price-story/pkg/router"
)

// methodSet is a tiny helper to store HTTP methods for a path
type methodSet map[string]struct{}

// specRoutes parses the OpenAPI YAML and returns a map of path -> methods
func specRoutes(t *testing.T, specPath string) map[string]methodSet {
	t.Helper()
	b, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("failed to read OpenAPI spec: %v", err)
	}
	var root map[string]any
	if err := yaml.Unmarshal(b, &root); err != nil {
		t.Fatalf("failed to unmarshal OpenAPI spec: %v", err)
	}
	pathsNode, ok := root["paths"].(map[string]any)
	if !ok {
		t.Fatalf("OpenAPI spec missing 'paths' section")
	}
	routes := make(map[string]methodSet)
	for p, v := range pathsNode {
		ops, ok := v.(map[string]any)
		if !ok {
			// tolerate non-object nodes
			continue
		}
		for _, m := range []string{"get", "post", "put", "delete"} {
			if op, exists := ops[m]; exists && op != nil {
				if _, ok := routes[p]; !ok {
					routes[p] = make(methodSet)
				}
				routes[p][strings.ToUpper(m)] = struct{}{}
			}
		}
	}
	return routes
}

// routerRoutes walks the Gorilla Mux router and returns a map of path -> methods
func routerRoutes(t *testing.T, r *sql.DB) map[string]methodSet {
	t.Helper()
	muxRouter := router.New(r)
	got := make(map[string]methodSet)
	walkErr := muxRouter.Walk(func(rt *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		tpl, err := rt.GetPathTemplate()
		if err != nil {
			return nil // skip routes without a template
		}
		methods, _ := rt.GetMethods()
		if len(methods) == 0 {
			return nil
		}
		if _, ok := got[tpl]; !ok {
			got[tpl] = make(methodSet)
		}
		for _, m := range methods {
			got[tpl][m] = struct{}{}
		}
		return nil
	})
	if walkErr != nil {
		t.Fatalf("router.Walk error: %v", walkErr)
	}
	return got
}

// prepareMockDB sets expectations so router.New can ensure the system actor user
func prepareMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	email := "system@example.com"
	pass := "placeholder"

	// EnsureUserByEmail first tries GetUserByEmail and then inserts with defaults when not found
	mock.ExpectQuery(regexp.QuoteMeta(queries.GetUserByEmail)).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery(regexp.QuoteMeta(queries.InsertUserReturning)).
		WithArgs(email, pass, "USD", 2, ",", "before").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "email", "password_hash", "preferred_currency", "decimal_places", "thousand_separator", "currency_symbol_placement",
		}).AddRow("sys", email, pass, "USD", 2, ",", "before"))

	return db, mock
}

func TestOpenAPIRoutesSync(t *testing.T) {
	// Resolve spec path relative to this package directory
	spec := filepath.Join("..", "..", "cmd", "price-story", "openapi", "openapi.yaml")

	// Build router with a mocked DB so startup middleware succeeds
	db, mock := prepareMockDB(t)
	defer func() { _ = db.Close() }()

	specMap := specRoutes(t, spec)
	routerMap := routerRoutes(t, db)

	// Compare: every router route must exist in spec
	for p, methods := range routerMap {
		sm, ok := specMap[p]
		if !ok {
			t.Errorf("router path %s is missing from OpenAPI spec", p)
			continue
		}
		for m := range methods {
			if _, ok := sm[m]; !ok {
				t.Errorf("router route %s %s missing from OpenAPI spec", m, p)
			}
		}
	}

	// Compare: every spec route must exist in router
	for p, methods := range specMap {
		rm, ok := routerMap[p]
		if !ok {
			t.Errorf("OpenAPI path %s is not registered in router", p)
			continue
		}
		for m := range methods {
			if _, ok := rm[m]; !ok {
				t.Errorf("OpenAPI route %s %s not registered in router", m, p)
			}
		}
	}

	// Ensure sqlmock expectations met (system user injection)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations not met: %v", err)
	}
}

// Additional validation: writable ops should define requestBody; id params must be documented
func TestOpenAPIOperationDetails(t *testing.T) {
	spec := filepath.Join("..", "..", "cmd", "price-story", "openapi", "openapi.yaml")
	b, err := os.ReadFile(spec)
	if err != nil {
		t.Fatalf("failed to read OpenAPI spec: %v", err)
	}
	var root map[string]any
	if err := yaml.Unmarshal(b, &root); err != nil {
		t.Fatalf("failed to unmarshal OpenAPI spec: %v", err)
	}
	pathsNode, ok := root["paths"].(map[string]any)
	if !ok {
		t.Fatalf("OpenAPI spec missing 'paths' section")
	}
	for p, v := range pathsNode {
		ops, ok := v.(map[string]any)
		if !ok {
			continue
		}
		for _, m := range []string{"post", "put"} {
			if opRaw, exists := ops[m]; exists && opRaw != nil {
				op, _ := opRaw.(map[string]any)
				if _, present := op["requestBody"]; !present {
					t.Errorf("OpenAPI operation %s %s missing requestBody", strings.ToUpper(m), p)
				}
			}
		}
		if strings.Contains(p, "{id}") {
			for _, m := range []string{"get", "put", "delete"} {
				if opRaw, exists := ops[m]; exists && opRaw != nil {
					op, _ := opRaw.(map[string]any)
					params, _ := op["parameters"].([]any)
					found := false
					for _, pr := range params {
						pm, _ := pr.(map[string]any)
						name, _ := pm["name"].(string)
						inVal, _ := pm["in"].(string)
						if name == "id" && inVal == "path" {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("OpenAPI operation %s %s missing path parameter 'id'", strings.ToUpper(m), p)
					}
				}
			}
		}
	}
}
