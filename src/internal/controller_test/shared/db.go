package shared

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
)

// LoadTestEnv loads .env.test.local from common project directories.
func LoadTestEnv() {
	envFiles := []string{
		os.Getenv("ENV_FILE"),
		".env.test.local",
		"../.env.test.local",
		"../../.env.test.local",
		"../../../.env.test.local",
		"/go/app/.env.test.local",
		"/home/thk/platform/platform-go/src/.env.test.local",
		"/home/thk/platform/platform-go/.env.test.local",
	}
	for _, f := range envFiles {
		if f != "" {
			if _, err := os.Stat(f); err == nil {
				_ = godotenv.Overload(f)
				return
			}
		}
	}
}

// SetupTestDB loads .env.test.local, connects to platform-db-test, ensures tables exist,
// and registers a cleanup function that wipes all data clean if the test succeeds.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	LoadTestEnv()

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "platform-db-test"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "3306"
		}
		user := os.Getenv("DB_USER")
		if user == "" {
			user = "root"
		}
		pass := os.Getenv("DB_PASS")
		if pass == "" {
			pass = "root"
		}
		dbname := os.Getenv("DB_NAME")
		if dbname == "" {
			dbname = "platform_db_test"
		}
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbname)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("skipping integration test: database unavailable: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil || sqlDB.Ping() != nil {
		t.Skip("skipping integration test: database unreachable")
	}

	// Ensure base tables exist if migrations were not executed beforehand
	migrator := db.Migrator()
	for _, model := range []any{
		&dbmodel.User{},
		&dbmodel.Homestay{},
		&dbmodel.HomestayRoom{},
		&dbmodel.HomestayImage{},
		&dbmodel.TravelLocation{},
	} {
		if !migrator.HasTable(model) {
			_ = db.AutoMigrate(model)
		}
	}

	// Clean up data on successful test run
	t.Cleanup(func() {
		if !t.Failed() {
			if err := CleanTestDB(db); err != nil {
				t.Logf("failed to clean test db: %v", err)
			}
		}
	})

	return db
}

// CleanTestDB clears all data from all tables in the current test database while preserving table schemas.
func CleanTestDB(db *gorm.DB) error {
	var tables []string
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE() AND table_type = 'BASE TABLE'"
	if err := db.Raw(query).Scan(&tables).Error; err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
			return err
		}
		for _, table := range tables {
			if table == "schema_migrations" {
				continue
			}
			if err := tx.Exec("TRUNCATE TABLE `" + table + "`").Error; err != nil {
				if err := tx.Exec("DELETE FROM `" + table + "`").Error; err != nil {
					return err
				}
			}
		}
		return tx.Exec("SET FOREIGN_KEY_CHECKS = 1").Error
	})
}
