package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/di"
	"platform-go/internal/infrastructure/security"
	"platform-go/internal/middleware"
	"platform-go/internal/pkg/app"
)

func main() {
	loadEnv()

	// Ưu tiên dùng DB_DSN nếu có, nếu không thì build từ từng biến lẻ
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASS")
		host := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")

		// Fallback mặc định nếu thiếu biến môi trường lẻ
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

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	jwtService, err := security.NewJWTService(jwtSecret)
	if err != nil {
		log.Fatalf("invalid JWT configuration: %v", err)
	}
	application := di.InitializeApplication(db, jwtService, !app.IsLocal())

	r := gin.Default()
	if err := r.SetTrustedProxies(trustedProxies()); err != nil {
		log.Fatalf("invalid TRUSTED_PROXIES configuration: %v", err)
	}

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		origin := os.Getenv("CORS_ALLOW_ORIGIN")
		if origin == "" {
			origin = "http://localhost:3000"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Vary", "Origin")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	r.Use(middleware.JWTAuth(application.AuthUsecase), middleware.CSRFProtection())

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
	if app.IsLocal() {
		// Nạp file .env.local
		err := godotenv.Overload(".env.local")
		if err != nil {
			log.Println("Không thể nạp file .env.local, sử dụng biến môi trường hệ thống")
		}
	} else if app.IsStaging() {
		// Nạp file .env.staging hoặc env.staging nếu có
		err := godotenv.Overload(".env.staging")
		if err != nil {
			err = godotenv.Overload("env.staging")
		}
		if err != nil {
			log.Println("Không thể nạp file .env.staging/env.staging, sử dụng biến môi trường hệ thống")
		}
	}
}
