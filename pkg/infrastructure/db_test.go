package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	file := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return file
}

func TestLoadConfig_FromYAML(t *testing.T) {
	yaml := `database:
  host: "db"
  port: 5432
  user: "user"
  password: "pass"
  dbname: "name"
  sslmode: "disable"
auth:
  jwt_secret: "secret"
  jwt_expiry_access_minutes: 15
  jwt_expiry_refresh_hours: 24
server:
  port: ":8080"
  allow_origins: "http://example.com"`
	file := writeTempConfig(t, yaml)

	cfg, err := LoadConfig(file)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.Database.Host != "db" || cfg.Database.Port != 5432 || cfg.Database.User != "user" ||
		cfg.Database.Password != "pass" || cfg.Database.DBName != "name" || cfg.Database.SSLMode != "disable" {
		t.Errorf("database values not loaded correctly: %+v", cfg.Database)
	}
	if cfg.Auth.JWTSecret != "secret" || cfg.Auth.JWTExpiryAccessMin != 15 || cfg.Auth.JWTExpiryRefreshHours != 24 {
		t.Errorf("auth values not loaded correctly: %+v", cfg.Auth)
	}
	if cfg.Server.Port != ":8080" {
		t.Errorf("server port not loaded correctly: %s", cfg.Server.Port)
	}
	if cfg.Server.AllowOrigins != "http://example.com" {
		t.Errorf("allow origins not loaded correctly: %s", cfg.Server.AllowOrigins)
	}
}

func TestLoadConfig_EnvOverride(t *testing.T) {
	yaml := `database:
  host: "db"
  port: 5432
  user: "user"
  password: "pass"
  dbname: "name"
  sslmode: "disable"
auth:
  jwt_secret: "secret"
  jwt_expiry_access_minutes: 15
  jwt_expiry_refresh_hours: 24
server:
  port: ":8080"
  allow_origins: "http://default.com"`
	file := writeTempConfig(t, yaml)

	t.Setenv("DB_HOST", "envhost")
	t.Setenv("DB_PORT", "6000")
	t.Setenv("DB_USER", "envuser")
	t.Setenv("DB_PASSWORD", "envpass")
	t.Setenv("DB_NAME", "envname")
	t.Setenv("DB_SSLMODE", "require")
	t.Setenv("JWT_SECRET", "envsecret")
	t.Setenv("JWT_EXPIRY_ACCESS", "30")
	t.Setenv("JWT_EXPIRY_REFRESH", "48")
	t.Setenv("SERVER_PORT", ":9999")
	t.Setenv("ALLOW_ORIGINS", "http://env.example.com")

	cfg, err := LoadConfig(file)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.Database.Host != "envhost" || cfg.Database.Port != 6000 || cfg.Database.User != "envuser" ||
		cfg.Database.Password != "envpass" || cfg.Database.DBName != "envname" || cfg.Database.SSLMode != "require" {
		t.Errorf("database env override failed: %+v", cfg.Database)
	}
	if cfg.Auth.JWTSecret != "envsecret" || cfg.Auth.JWTExpiryAccessMin != 30 || cfg.Auth.JWTExpiryRefreshHours != 48 {
		t.Errorf("auth env override failed: %+v", cfg.Auth)
	}
	if cfg.Server.Port != ":9999" {
		t.Errorf("server port env override failed: %s", cfg.Server.Port)
	}
	if cfg.Server.AllowOrigins != "http://env.example.com" {
		t.Errorf("allow origins env override failed: %s", cfg.Server.AllowOrigins)
	}
}

func TestLoadConfig_JWTSecretFile(t *testing.T) {
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	if err := os.WriteFile(secretFile, []byte("file-secret"), 0o644); err != nil {
		t.Fatalf("failed to write secret file: %v", err)
	}
	yaml := `database:
  host: "db"
  port: 5432
  user: "user"
  password: "pass"
  dbname: "name"
  sslmode: "disable"
auth:
  jwt_secret: "%s"
  jwt_expiry_access_minutes: 15
  jwt_expiry_refresh_hours: 24
server:
  port: ":8080"
  allow_origins: "http://file.com"`
	yaml = fmt.Sprintf(yaml, secretFile)
	file := writeTempConfig(t, yaml)

	cfg, err := LoadConfig(file)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if cfg.Auth.JWTSecret != "file-secret" {
		t.Errorf("expected secret from file, got %s", cfg.Auth.JWTSecret)
	}
	if cfg.Server.AllowOrigins != "http://file.com" {
		t.Errorf("allow origins not loaded correctly from file: %s", cfg.Server.AllowOrigins)
	}
}
