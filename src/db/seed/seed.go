package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"platform-go/db/seed/seeder"
	"platform-go/internal/infrastructure/repository"
)

func main() {
	// Load env from the corresponding file based on environment
	if os.Getenv("APP_ENV") == "staging" {
		if err := godotenv.Load(".env.staging"); err != nil {
			_ = godotenv.Load("env.staging")
		}
	} else {
		_ = godotenv.Load(".env.local")
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASS")
		host := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")

		// Fallback defaults if running outside Docker or missing individual variables
		if user == "" {
			user = "root"
		}
		if pass == "" {
			pass = "root"
		}
		if host == "" {
			host = "localhost"
		} // Point to localhost if running from host
		if dbPort == "" {
			dbPort = "3309"
		}
		if dbName == "" {
			dbName = "platform_db"
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, dbPort, dbName)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	schemaDB := repository.NewSchemaDB(db)

	seeder.CreateUsers(schemaDB)

	//Homestay
	homestays := seeder.CreateHomestays(schemaDB)
	rooms := seeder.CreateHomestayRooms(schemaDB, homestays)
	seeder.CreateHomestayImages(schemaDB, rooms)

	//Travel map
	locations := seeder.CreateTravelLocations(schemaDB)

	//Travel memories
	memories := seeder.CreateTravelMemories(schemaDB, locations)
	days := seeder.CreateTravelMemoryDays(schemaDB, memories)
	seeder.CreateTravelMemoryPhotos(schemaDB, days)

	fmt.Println("Seeding completed successfully.")
}
