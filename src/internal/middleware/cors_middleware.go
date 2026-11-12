package middleware

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	corsAllowOriginsEnv = "CORS_ALLOW_ORIGIN"

	corsAllowHeaders = "Content-Type, Content-Length, Accept-Encoding, " +
		"X-CSRF-Token, Authorization, Accept, Origin, " +
		"Cache-Control, X-Requested-With"

	corsAllowMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
)

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS).
// Allowed origins are loaded and validated once when the middleware is initialized.
func CORSMiddleware() gin.HandlerFunc {
	rawOrigins := os.Getenv(corsAllowOriginsEnv)
	allowAll := false

	parts := strings.FieldsFunc(rawOrigins, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' ||
			r == '\n' || r == '\r'
	})

	for _, part := range parts {
		if strings.TrimSpace(part) == "*" {
			allowAll = true
			break
		}
	}

	allowedOrigins := parseAllowedOrigins(rawOrigins)

	allowedOriginSet := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowedOriginSet[origin] = struct{}{}
	}

	if !allowAll && len(allowedOriginSet) == 0 {
		log.Printf(
			"[WARN] %s is missing or contains no valid origins; cross-origin requests will not be allowed",
			corsAllowOriginsEnv,
		)
	}

	return func(c *gin.Context) {
		header := c.Writer.Header()
		addVaryOrigin(header)

		requestOrigin := c.GetHeader("Origin")

		if requestOrigin != "" {
			normalizedOrigin, valid := normalizeAllowedOrigin(requestOrigin)

			if valid {
				_, allowed := allowedOriginSet[normalizedOrigin]
				if allowAll || allowed {
					allowHeaders := corsAllowHeaders
					if reqHeaders := c.GetHeader("Access-Control-Request-Headers"); reqHeaders != "" {
						allowHeaders = corsAllowHeaders + ", " + reqHeaders
					}

					header.Set("Access-Control-Allow-Origin", requestOrigin)
					header.Set("Access-Control-Allow-Credentials", "true")
					header.Set("Access-Control-Allow-Headers", allowHeaders)
					header.Set("Access-Control-Allow-Methods", corsAllowMethods)
					header.Set("Access-Control-Max-Age", "600")
				}
			}
		}

		// Handle preflight before authentication and CSRF middleware.
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// parseAllowedOrigins parses comma- or whitespace-separated origins.
func parseAllowedOrigins(raw string) []string {
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' ||
			r == '\n' || r == '\r'
	})

	origins := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))

	for _, part := range parts {
		if strings.TrimSpace(part) == "*" {
			continue
		}

		origin, ok := normalizeAllowedOrigin(part)
		if !ok {
			log.Printf("[WARN] Ignoring invalid CORS origin: %q", part)
			continue
		}

		if _, exists := seen[origin]; exists {
			continue
		}

		seen[origin] = struct{}{}
		origins = append(origins, origin)
	}

	return origins
}

// normalizeAllowedOrigin accepts HTTP(S) origins without paths,
// queries, fragments, credentials, or wildcard hosts.
// A trailing slash and default HTTP(S) ports are normalized away.
func normalizeAllowedOrigin(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)

	if raw == "" ||
		raw == "*" ||
		strings.EqualFold(raw, "null") ||
		strings.ContainsAny(raw, "* \t\r\n") {
		return "", false
	}

	parsed, err := url.Parse(raw)
	if err != nil ||
		parsed.Host == "" ||
		parsed.User != nil ||
		parsed.Opaque != "" ||
		(parsed.Path != "" && parsed.Path != "/") ||
		parsed.RawQuery != "" ||
		parsed.ForceQuery ||
		parsed.Fragment != "" {
		return "", false
	}

	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", false
	}

	hostname := strings.ToLower(parsed.Hostname())
	if hostname == "" {
		return "", false
	}

	// Reject malformed hostnames and invalid ports.
	port := parsed.Port()
	if strings.Contains(parsed.Host, ":") && port == "" {
		// A colon may be valid in an IPv6 literal, so validate below.
		if !strings.Contains(hostname, ":") {
			return "", false
		}
	}

	if port != "" {
		for _, char := range port {
			if char < '0' || char > '9' {
				return "", false
			}
		}
	}

	if (scheme == "http" && port == "80") ||
		(scheme == "https" && port == "443") {
		port = ""
	}

	host := hostname
	if strings.Contains(hostname, ":") {
		host = "[" + hostname + "]"
	}

	if port != "" {
		host += ":" + port
	}

	return scheme + "://" + host, true
}

// addVaryOrigin preserves existing Vary headers and adds Origin when needed.
func addVaryOrigin(header http.Header) {
	for _, value := range header.Values("Vary") {
		for _, token := range strings.Split(value, ",") {
			token = strings.TrimSpace(token)

			if token == "*" || strings.EqualFold(token, "Origin") {
				return
			}
		}
	}

	header.Add("Vary", "Origin")
}
