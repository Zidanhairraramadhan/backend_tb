package config

import (
	"log"
	"os"
	"strings"

	"musiclink-backend/model"

	"github.com/glebarez/sqlite" // Pure Go SQLite driver (no CGO required)
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

		// MODIFIKASI DI SINI: Menggunakan postgres.New untuk menambahkan PreferSimpleProtocol
		DB, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // Menonaktifkan prepared statements untuk kompatibilitas pooler Supabase (port 6543)
		}), &gorm.Config{})

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
	err = DB.AutoMigrate(&model.User{}, &model.Link{}, &model.ClickLog{})
	if err != nil {
		log.Fatalf("❌ Database Migration Failed: %v", err)
	}
	log.Println("✅ Database Migration completed successfully.")
}
