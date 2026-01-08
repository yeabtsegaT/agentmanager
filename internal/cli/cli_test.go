package cli

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"zero", 0, "0s"},
		{"seconds", 45 * time.Second, "45s"},
		{"one minute", time.Minute, "1m 0s"},
		{"minutes and seconds", 5*time.Minute + 30*time.Second, "5m 30s"},
		{"one hour", time.Hour, "1h 0m"},
		{"hours and minutes", 3*time.Hour + 15*time.Minute, "3h 15m"},
		{"one day", 24 * time.Hour, "1d 0h"},
		{"days and hours", 2*24*time.Hour + 5*time.Hour, "2d 5h"},
		{"many days", 10*24*time.Hour + 12*time.Hour, "10d 12h"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, result, tt.expected)
			}
		})
	}
}

func TestParseConfigValue(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		expected interface{}
	}{
		// Boolean keys
		{"bool true lowercase", "catalog.refresh_on_start", "true", true},
		{"bool true uppercase", "catalog.refresh_on_start", "TRUE", true},
		{"bool true yes", "updates.auto_check", "yes", true},
		{"bool true 1", "updates.notify", "1", true},
		{"bool false lowercase", "updates.auto_update", "false", false},
		{"bool false no", "ui.show_hidden", "no", false},
		{"bool false 0", "ui.use_colors", "0", false},
		{"bool false random", "ui.compact_mode", "random", false},
		{"bool api enable grpc", "api.enable_grpc", "true", true},
		{"bool api enable rest", "api.enable_rest", "false", false},
		{"bool api require auth", "api.require_auth", "yes", true},

		// Integer keys
		{"int ui page size", "ui.page_size", "50", 50},
		{"int grpc port", "api.grpc_port", "50051", 50051},
		{"int rest port", "api.rest_port", "8080", 8080},
		{"int logging max size", "logging.max_size", "100", 100},
		{"int logging max age", "logging.max_age", "30", 30},
		{"int invalid returns string", "ui.page_size", "invalid", "invalid"},

		// Duration keys
		{"duration refresh interval", "catalog.refresh_interval", "1h", time.Hour},
		{"duration check interval", "updates.check_interval", "30m", 30 * time.Minute},
		{"duration with seconds", "catalog.refresh_interval", "1h30m45s", time.Hour + 30*time.Minute + 45*time.Second},
		{"duration invalid returns string", "catalog.refresh_interval", "invalid", "invalid"},

		// String keys (default)
		{"string unknown key", "unknown.key", "some value", "some value"},
		{"string catalog url", "catalog.source_url", "https://example.com", "https://example.com"},
		{"string with spaces", "some.key", "value with spaces", "value with spaces"},

		// Case insensitive keys
		{"case insensitive bool", "CATALOG.REFRESH_ON_START", "true", true},
		{"case insensitive int", "UI.PAGE_SIZE", "25", 25},
		{"case insensitive duration", "CATALOG.REFRESH_INTERVAL", "2h", 2 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseConfigValue(tt.key, tt.value)
			if result != tt.expected {
				t.Errorf("parseConfigValue(%q, %q) = %v (%T), want %v (%T)",
					tt.key, tt.value, result, result, tt.expected, tt.expected)
			}
		})
	}
}

func TestParseConfigValueAllBoolKeys(t *testing.T) {
	boolKeys := []string{
		"catalog.refresh_on_start",
		"updates.auto_check",
		"updates.notify",
		"updates.auto_update",
		"ui.show_hidden",
		"ui.use_colors",
		"ui.compact_mode",
		"api.enable_grpc",
		"api.enable_rest",
		"api.require_auth",
	}

	for _, key := range boolKeys {
		t.Run(key, func(t *testing.T) {
			// Test that "true" returns true
			result := parseConfigValue(key, "true")
			if result != true {
				t.Errorf("parseConfigValue(%q, %q) = %v, want true", key, "true", result)
			}

			// Test that "false" returns false
			result = parseConfigValue(key, "false")
			if result != false {
				t.Errorf("parseConfigValue(%q, %q) = %v, want false", key, "false", result)
			}
		})
	}
}

func TestParseConfigValueAllIntKeys(t *testing.T) {
	intKeys := []string{
		"ui.page_size",
		"api.grpc_port",
		"api.rest_port",
		"logging.max_size",
		"logging.max_age",
	}

	for _, key := range intKeys {
		t.Run(key, func(t *testing.T) {
			// Test that a valid integer is parsed
			result := parseConfigValue(key, "42")
			if result != 42 {
				t.Errorf("parseConfigValue(%q, %q) = %v (%T), want 42 (int)", key, "42", result, result)
			}
		})
	}
}

func TestParseConfigValueAllDurationKeys(t *testing.T) {
	durationKeys := []string{
		"catalog.refresh_interval",
		"updates.check_interval",
	}

	for _, key := range durationKeys {
		t.Run(key, func(t *testing.T) {
			// Test that a valid duration is parsed
			result := parseConfigValue(key, "1h")
			if result != time.Hour {
				t.Errorf("parseConfigValue(%q, %q) = %v (%T), want 1h (time.Duration)", key, "1h", result, result)
			}
		})
	}
}
