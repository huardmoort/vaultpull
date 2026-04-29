package env_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/env"
)

func TestIsSensitive_CommonKeys(t *testing.T) {
	sensitiveKeys := []string{
		"PASSWORD",
		"DB_PASSWORD",
		"SECRET",
		"API_SECRET",
		"TOKEN",
		"AUTH_TOKEN",
		"PRIVATE_KEY",
		"AWS_SECRET_ACCESS_KEY",
		"CREDENTIALS",
		"APP_CREDENTIALS",
	}

	for _, key := range sensitiveKeys {
		t.Run(key, func(t *testing.T) {
			if !env.IsSensitive(key) {
				t.Errorf("expected IsSensitive(%q) = true, got false", key)
			}
		})
	}
}

func TestIsSensitive_SafeKeys(t *testing.T) {
	safeKeys := []string{
		"APP_NAME",
		"PORT",
		"LOG_LEVEL",
		"DATABASE_HOST",
		"REGION",
		"ENVIRONMENT",
	}

	for _, key := range safeKeys {
		t.Run(key, func(t *testing.T) {
			if env.IsSensitive(key) {
				t.Errorf("expected IsSensitive(%q) = false, got true", key)
			}
		})
	}
}

func TestRedactValue_SensitiveKey(t *testing.T) {
	result := env.RedactValue("DB_PASSWORD", "supersecret")
	if result == "supersecret" {
		t.Error("expected value to be redacted, but got original value")
	}
	if result == "" {
		t.Error("expected non-empty redacted placeholder")
	}
}

func TestRedactValue_SafeKey(t *testing.T) {
	result := env.RedactValue("APP_NAME", "myapp")
	if result != "myapp" {
		t.Errorf("expected value to be unchanged, got %q", result)
	}
}

func TestRedactValue_EmptyValue(t *testing.T) {
	// Even sensitive keys with empty values should return a consistent placeholder
	result := env.RedactValue("PASSWORD", "")
	if result == "" {
		t.Error("expected non-empty redacted placeholder for sensitive key with empty value")
	}
}

func TestRedactMap_MixedKeys(t *testing.T) {
	input := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "s3cr3t",
		"PORT":        "8080",
		"API_TOKEN":   "tok_abc123",
	}

	result := env.RedactMap(input)

	if result["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should be unchanged, got %q", result["APP_NAME"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("PORT should be unchanged, got %q", result["PORT"])
	}
	if result["DB_PASSWORD"] == "s3cr3t" {
		t.Error("DB_PASSWORD should be redacted")
	}
	if result["API_TOKEN"] == "tok_abc123" {
		t.Error("API_TOKEN should be redacted")
	}
}

func TestRedactMap_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{
		"DB_PASSWORD": "original",
	}

	_ = env.RedactMap(input)

	if input["DB_PASSWORD"] != "original" {
		t.Error("RedactMap must not mutate the input map")
	}
}

func TestRedactMap_EmptyInput(t *testing.T) {
	result := env.RedactMap(map[string]string{})
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}
