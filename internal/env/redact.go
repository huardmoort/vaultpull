// Package env provides utilities for validating and handling environment
// variable key-value pairs before they are written to .env files.
package env

import (
	"strings"
)

// sensitiveKeyPatterns holds substrings that indicate a key likely contains
// a sensitive value that should be redacted in logs or terminal output.
var sensitiveKeyPatterns = []string{
	"PASSWORD",
	"PASSWD",
	"SECRET",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE_KEY",
	"PRIVATE",
	"CREDENTIAL",
	"AUTH",
	"ACCESS_KEY",
	"CERT",
	"DSN",
	"DATABASE_URL",
	"DB_URL",
}

// redactedPlaceholder is the string substituted for sensitive values.
const redactedPlaceholder = "[REDACTED]"

// IsSensitive reports whether the given environment variable key is considered
// sensitive based on well-known naming conventions. The comparison is
// case-insensitive.
func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range sensitiveKeyPatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

// RedactValue returns the original value if the key is not considered
// sensitive, or redactedPlaceholder if it is. Use this when displaying
// key-value pairs in logs, diffs, or terminal output.
func RedactValue(key, value string) string {
	if IsSensitive(key) {
		return redactedPlaceholder
	}
	return value
}

// RedactMap returns a copy of the provided map with sensitive values replaced
// by redactedPlaceholder. The original map is not modified.
func RedactMap(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = RedactValue(k, v)
	}
	return out
}
