package config

import (
	"log"
	"os"
	"strings"

	"github.com/glebarez/sqlite" // Pure Go SQLite driver (no CGO required)
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"musiclink-backend/model"
)

var DB *gorm.DB

func maskDSN(dsn string) string {
	if dsn == "" {
		return "[EMPTY]"
	}
	parts := strings.Split(dsn, "@")
	if len(parts) > 1 {
		left := parts[0]
		passIndex := strings.LastIndex(left, ":")
		if passIndex != -1 {
			return left[:passIndex] + ":****@" + parts[1]
		}
	}
	return dsn
}

func ConnectDB() {
	var err error
	dsn := os.Getenv("SUPABASE_DSN")
	log.Printf("🔌 Loaded DSN: %s", maskDSN(dsn))

	if dsn != "" {
		log.Println("🔌 Connecting to Supabase PostgreSQL...")
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("⚠️ Failed to connect to Supabase PostgreSQL: %v\n", err)
			log.Println("🔄 Falling back to local pure-Go SQLite database (gorm.db)...")
			DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
		}
	} else {
		log.Println("ℹ️ SUPABASE_DSN not set or empty. Connecting to local pure-Go SQLite database (gorm.db)...")
		DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	}

	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	log.Println("✅ Database connection established.")

	// Auto Migration
	log.Println("🔧 Running GORM AutoMigrations...")
	err = DB.AutoMigrate(&model.User{}, &model.Link{})
	if err != nil {
		log.Fatalf("❌ Database Migration Failed: %v", err)
	}
	log.Println("✅ Database Migration completed successfully.")
}
