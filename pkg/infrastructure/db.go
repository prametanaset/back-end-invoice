package infrastructure

import (
	"fmt"
	"log"
	"os"
	"strings"

	"io/ioutil"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBConfig คือ struct เก็บค่าตั้งค่าฐานข้อมูล
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

// AppConfig เก็บทั้งส่วน Database และ Auth
type AppConfig struct {
	Database DBConfig `yaml:"database"`
	Auth     struct {
		JWTSecret             string `yaml:"jwt_secret"`
		JWTExpiryAccessMin    int    `yaml:"jwt_expiry_access_minutes"`
		JWTExpiryRefreshHours int    `yaml:"jwt_expiry_refresh_hours"`
	} `yaml:"auth"`
	Server ServerConfig `yaml:"server"`
}

// LoadConfig อ่าน config จากไฟล์ YAML แล้ว merge กับ ENV ได้
func LoadConfig(path string) (*AppConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	// ถ้ามี ENV กำหนดมาให้ override ค่าใน cfg ได้
	// ตัวอย่างเช่น ENV: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE
	// หรือ ENV: JWT_SECRET, JWT_EXPIRY_ACCESS, JWT_EXPIRY_REFRESH
	if env := os.Getenv("DB_HOST"); env != "" {
		cfg.Database.Host = env
	}
	if env := os.Getenv("DB_PORT"); env != "" {
		// แปลง string เป็น int
		fmt.Sscanf(env, "%d", &cfg.Database.Port)
	}
	if env := os.Getenv("DB_USER"); env != "" {
		cfg.Database.User = env
	}
	if env := os.Getenv("DB_PASSWORD"); env != "" {
		cfg.Database.Password = env
	}
	if env := os.Getenv("DB_NAME"); env != "" {
		cfg.Database.DBName = env
	}
	if env := os.Getenv("DB_SSLMODE"); env != "" {
		cfg.Database.SSLMode = env
	}
	// สำหรับ JWT secret
	if env := os.Getenv("JWT_SECRET"); env != "" {
		cfg.Auth.JWTSecret = env
	}
	if env := os.Getenv("JWT_EXPIRY_ACCESS"); env != "" {
		fmt.Sscanf(env, "%d", &cfg.Auth.JWTExpiryAccessMin)
	}
	if env := os.Getenv("JWT_EXPIRY_REFRESH"); env != "" {
		fmt.Sscanf(env, "%d", &cfg.Auth.JWTExpiryRefreshHours)
	}
	// If JWTSecret points to a file, read its contents
	if cfg.Auth.JWTSecret != "" {
		if b, err := os.ReadFile(cfg.Auth.JWTSecret); err == nil {
			cfg.Auth.JWTSecret = strings.TrimSpace(string(b))
		}
	}
	return &cfg, nil
}

// NewPostgresDB สร้าง *gorm.DB จากค่าที่ได้
func NewPostgresDB(dbCfg DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.DBName, dbCfg.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Migrate ทำ AutoMigrate ให้กับโมเดลต่าง ๆ (เพิ่มตาราง RefreshToken ด้วย)
func Migrate(db *gorm.DB, models ...interface{}) {
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
}
