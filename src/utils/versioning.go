package utils

import (
	"context"
	"net/http"
	"regexp"

	"github.com/fatih/color"
)

// Key type for storing the API version in the request context
type contextKey string

const (
	ApiVersionKey contextKey = "apiVersion"
)

// VersionMiddleware is a custom middleware for extracting the API version from the header
func VersionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiVersion := r.Header.Get("X-Version")

		// Check if the version is allowed
		if !isValidVersion(apiVersion) {
			// Handle unsupported version
			color.Red("Unsupported api version type", apiVersion)
			apiVersion = ""
		}
		// Set the version in the request context
		ctx := context.WithValue(r.Context(), ApiVersionKey, apiVersion)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// isValidVersion checks if the provided version string is valid
func isValidVersion(version string) bool {
	if version == "" {
		return true
	}

	// Define a regular expression pattern for the allowed version format
	allowedPattern := "^[0-9.%]+$"
	match, err := regexp.MatchString(allowedPattern, version)
	if err != nil {
		color.Red("error during matching", err)
		return false
	}
	return match
}
