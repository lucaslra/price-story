package config

import (
    "testing"
    "time"
)

func TestNew_Defaults(t *testing.T) {
    cfg := New()
    if cfg.Server.Port != "8080" {
        t.Fatalf("expected default port 8080, got %s", cfg.Server.Port)
    }
    if cfg.Server.ReadTimeout != 15*time.Second || cfg.Server.WriteTimeout != 15*time.Second {
        t.Fatalf("expected default timeouts 15s, got read=%v write=%v", cfg.Server.ReadTimeout, cfg.Server.WriteTimeout)
    }
    if cfg.Database.Host != "localhost" || cfg.Database.Port != 5432 || cfg.Database.User != "postgres" || cfg.Database.Password != "postgres" || cfg.Database.DBName != "price_story" || cfg.Database.SSLMode != "disable" {
        t.Fatalf("unexpected default db config: %+v", cfg.Database)
    }
}

func TestNew_EnvOverrides(t *testing.T) {
    t.Setenv("PORT", "9000")
    t.Setenv("READ_TIMEOUT", "20s")
    t.Setenv("WRITE_TIMEOUT", "25s")
    t.Setenv("DB_HOST", "db.local")
    t.Setenv("DB_PORT", "6543")
    t.Setenv("DB_USER", "user")
    t.Setenv("DB_PASSWORD", "pass")
    t.Setenv("DB_NAME", "db")
    t.Setenv("DB_SSLMODE", "require")

    cfg := New()
    if cfg.Server.Port != "9000" {
        t.Fatalf("expected port 9000, got %s", cfg.Server.Port)
    }
    if cfg.Server.ReadTimeout != 20*time.Second || cfg.Server.WriteTimeout != 25*time.Second {
        t.Fatalf("expected timeouts 20s/25s, got read=%v write=%v", cfg.Server.ReadTimeout, cfg.Server.WriteTimeout)
    }
    if cfg.Database.Host != "db.local" || cfg.Database.Port != 6543 || cfg.Database.User != "user" || cfg.Database.Password != "pass" || cfg.Database.DBName != "db" || cfg.Database.SSLMode != "require" {
        t.Fatalf("unexpected overridden db config: %+v", cfg.Database)
    }
}