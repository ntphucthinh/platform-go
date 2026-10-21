package app

import "os"

// IsLocal determines if the application is running in a local environment
func IsLocal() bool {
	env := os.Getenv("APP_ENV")
	return env == "" || env == "local"
}

// IsStaging determines if the application is running in a staging environment
func IsStaging() bool {
	return os.Getenv("APP_ENV") == "staging"
}

// Environment returns the current application environment
func Environment() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		return "local"
	}
	return env
}

