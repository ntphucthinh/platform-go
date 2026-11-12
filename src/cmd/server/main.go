package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	api "platform-go/generated/api"
	"platform-go/internal/di"
	"platform-go/internal/infrastructure/security"
	"platform-go/internal/middleware"
	"platform-go/internal/pkg/app"
	"platform-go/internal/pkg/cvalidator"
)

func main() {
	loadEnv()

	// Prefer DB_DSN if set; otherwise build from individual environment variables
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASS")
		host := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")

		// Fallback defaults if individual variables are missing
		if user == "" {
			user = "root"
		}
		if pass == "" {
			pass = "root"
		}
		if host == "" {
			host = "platform-db"
		}
		if dbPort == "" {
			dbPort = "3306"
		}
		if dbName == "" {
			dbName = "platform_db"
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, dbPort, dbName)
	}

	logLevel := logger.Warn
	if app.IsLocal() {
		logLevel = logger.Info
	}
	switch strings.ToLower(os.Getenv("DB_LOG_LEVEL")) {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	}

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	jwtService, err := security.NewJWTService(jwtSecret)
	if err != nil {
		log.Fatalf("invalid JWT configuration: %v", err)
	}
	application := di.InitializeApplication(db, jwtService, !app.IsLocal())

	cvalidator.NewValidate()
	v := cvalidator.Validator()
	cvalidator.CustomTagName(v)
	cvalidator.ValidatorTranslate(v)
	cvalidator.AddCustomValidation(v, application.DB)

	r := gin.Default()
	if err := r.SetTrustedProxies(trustedProxies()); err != nil {
		log.Fatalf("invalid TRUSTED_PROXIES configuration: %v", err)
	}

	// CORS Middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.JWTAuth(application.AuthUsecase))
	r.Use(middleware.CSRFProtection())

	r.GET("/health", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	// Register generated routes to the gin router
	api.RegisterHandlers(r, application.Controller)
	// Keep the former unprefixed endpoints available while clients migrate to /auth/*.
	r.POST("/login", application.Controller.Login)
	r.POST("/logout", application.Controller.Logout)
	r.GET("/me", application.Controller.Me)
	r.POST("/refresh", application.Controller.Refresh)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// trustedProxies returns no trusted proxies by default. Configure explicit proxy
// IPs/CIDRs only when the application is deployed behind a trusted proxy.
func trustedProxies() []string {
	raw := strings.TrimSpace(os.Getenv("TRUSTED_PROXIES"))
	if raw == "" {
		return nil
	}
	return strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' || r == '\n'
	})
}

func loadEnv() {
	// 1. Cho phép chỉ định file cấu hình cụ thể qua biến ENV_FILE nếu cần
	if envFile := os.Getenv("ENV_FILE"); envFile != "" {
		if err := godotenv.Overload(envFile); err == nil {
			log.Printf("[INFO] Loaded environment from %s", envFile)
			return
		}
		log.Printf("[WARN] Could not load specified ENV_FILE %q, falling back", envFile)
	}

	// 2. Nếu đã khai báo rõ APP_ENV từ shell/system (vd: export APP_ENV=staging từ service.sh hoặc Docker)
	appEnv := strings.TrimSpace(os.Getenv("APP_ENV"))
	if appEnv != "" && appEnv != "local" {
		targetFile := ".env." + appEnv
		if _, err := os.Stat(targetFile); err == nil {
			if err := godotenv.Overload(targetFile); err == nil {
				log.Printf("[INFO] Loaded environment from %s", targetFile)
				return
			}
		}
		if _, err := os.Stat(".env"); err == nil {
			if err := godotenv.Overload(".env"); err == nil {
				log.Printf("[INFO] Loaded environment from .env (for APP_ENV=%s)", appEnv)
				return
			}
		}
	}

	// 3. Môi trường local dev: ưu tiên .env.local, sau đó đến .env
	if _, err := os.Stat(".env.local"); err == nil {
		if err := godotenv.Overload(".env.local"); err == nil {
			log.Println("[INFO] Loaded environment from .env.local")
			return
		}
		log.Println("[WARN] Found .env.local but could not parse it, falling back")
	}

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Overload(".env"); err == nil {
			log.Println("[INFO] Loaded environment from .env")
			return
		}
		log.Println("[WARN] Found .env but could not parse it, falling back")
	}

	// 4. Dự phòng nếu chỉ có .env.staging hoặc .env.production
	if _, err := os.Stat(".env.staging"); err == nil {
		_ = godotenv.Overload(".env.staging")
		log.Println("[INFO] Loaded environment from .env.staging")
	} else if _, err := os.Stat(".env.production"); err == nil {
		_ = godotenv.Overload(".env.production")
		log.Println("[INFO] Loaded environment from .env.production")
	}
}
