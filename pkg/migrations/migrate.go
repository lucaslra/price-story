package migrations

import (
	"database/sql"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Migration represents a database migration
type Migration struct {
	Version int64
	Name    string
	UpSQL   string
	DownSQL string
}

// Migrator handles database migrations
type Migrator struct {
	db           *sql.DB
	migrations   []Migration
	migrationsFS fs.FS
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB, migrationsFS fs.FS) *Migrator {
	return &Migrator{
		db:           db,
		migrationsFS: migrationsFS,
	}
}

// LoadMigrations loads all migration files from the filesystem
func (m *Migrator) LoadMigrations() error {
	migrations := make(map[int64]*Migration)

	err := fs.WalkDir(m.migrationsFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		// Parse filename: 001_create_table.up.sql or 001_create_table.down.sql
		filename := filepath.Base(path)
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			return fmt.Errorf("invalid migration filename: %s", filename)
		}

		versionStr := parts[0]
		version, err := strconv.ParseInt(versionStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid version in filename %s: %w", filename, err)
		}

		// Read file content
		content, err := fs.ReadFile(m.migrationsFS, path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", path, err)
		}

		// Get or create migration
		migration, exists := migrations[version]
		if !exists {
			migration = &Migration{
				Version: version,
				Name:    strings.Join(parts[1:len(parts)-1], "_"), // Remove version and .up/.down.sql
			}
			migrations[version] = migration
		}

		// Set SQL content based on type
		if strings.Contains(filename, ".up.sql") {
			migration.UpSQL = string(content)
		} else if strings.Contains(filename, ".down.sql") {
			migration.DownSQL = string(content)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Convert map to sorted slice
	m.migrations = make([]Migration, 0, len(migrations))
	for _, migration := range migrations {
		m.migrations = append(m.migrations, *migration)
	}

	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	return nil
}

// createMigrationsTable creates the schema_migrations table if it doesn't exist
func (m *Migrator) createMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version BIGINT PRIMARY KEY,
		dirty BOOLEAN NOT NULL DEFAULT FALSE,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	return nil
}

// getAppliedMigrations returns a map of applied migration versions
func (m *Migrator) getAppliedMigrations() (map[int64]bool, error) {
	applied := make(map[int64]bool)

	rows, err := m.db.Query("SELECT version FROM schema_migrations WHERE NOT dirty")
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[version] = true
	}

	return applied, nil
}

// Up runs all pending migrations
func (m *Migrator) Up() error {
	if err := m.createMigrationsTable(); err != nil {
		return err
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	for _, migration := range m.migrations {
		if applied[migration.Version] {
			fmt.Printf("Migration %d (%s) already applied, skipping\n", migration.Version, migration.Name)
			continue
		}

		if migration.UpSQL == "" {
			fmt.Printf("Skipping migration %d (%s): no up SQL\n", migration.Version, migration.Name)
			continue
		}

		fmt.Printf("Applying migration %d (%s)...\n", migration.Version, migration.Name)

		// Start transaction
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction for migration %d: %w", migration.Version, err)
		}

		// Mark as dirty
		_, err = tx.Exec("INSERT INTO schema_migrations (version, dirty) VALUES ($1, TRUE)", migration.Version)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to mark migration %d as dirty: %w", migration.Version, err)
		}

		// Execute migration
		_, err = tx.Exec(migration.UpSQL)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to execute migration %d: %w", migration.Version, err)
		}

		// Mark as clean
		_, err = tx.Exec("UPDATE schema_migrations SET dirty = FALSE WHERE version = $1", migration.Version)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to mark migration %d as clean: %w", migration.Version, err)
		}

		// Commit transaction
		if err = tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
		}

		fmt.Printf("Migration %d (%s) applied successfully\n", migration.Version, migration.Name)
	}

	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down() error {
	if err := m.createMigrationsTable(); err != nil {
		return err
	}

	// Get the latest applied migration
	var latestVersion int64
	err := m.db.QueryRow("SELECT version FROM schema_migrations WHERE NOT dirty ORDER BY version DESC LIMIT 1").Scan(&latestVersion)
	if err == sql.ErrNoRows {
		fmt.Println("No migrations to roll back")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get latest migration: %w", err)
	}

	// Find the migration
	var targetMigration *Migration
	for _, migration := range m.migrations {
		if migration.Version == latestVersion {
			targetMigration = &migration
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("migration %d not found in migration files", latestVersion)
	}

	if targetMigration.DownSQL == "" {
		return fmt.Errorf("migration %d (%s) has no down SQL", targetMigration.Version, targetMigration.Name)
	}

	fmt.Printf("Rolling back migration %d (%s)...\n", targetMigration.Version, targetMigration.Name)

	// Start transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction for rollback %d: %w", targetMigration.Version, err)
	}

	// Execute rollback
	_, err = tx.Exec(targetMigration.DownSQL)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to execute rollback %d: %w", targetMigration.Version, err)
	}

	// Remove from migrations table
	_, err = tx.Exec("DELETE FROM schema_migrations WHERE version = $1", targetMigration.Version)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to remove migration %d from schema_migrations: %w", targetMigration.Version, err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback %d: %w", targetMigration.Version, err)
	}

	fmt.Printf("Migration %d (%s) rolled back successfully\n", targetMigration.Version, targetMigration.Name)
	return nil
}

// Status shows the current migration status
func (m *Migrator) Status() error {
	if err := m.createMigrationsTable(); err != nil {
		return err
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")

	for _, migration := range m.migrations {
		status := "PENDING"
		if applied[migration.Version] {
			status = "APPLIED"
		}
		fmt.Printf("%d\t%s\t%s\n", migration.Version, status, migration.Name)
	}

	return nil
}
