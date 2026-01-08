package strategies

import (
	"context"
	"testing"

	"github.com/kevinelliott/agentmgr/pkg/agent"
	"github.com/kevinelliott/agentmgr/pkg/platform"
)

// mockPlatform implements platform.Platform for testing
type mockPlatform struct {
	executables map[string]string
}

func newMockPlatform() *mockPlatform {
	return &mockPlatform{
		executables: make(map[string]string),
	}
}

func (m *mockPlatform) ID() platform.ID                                             { return platform.Darwin }
func (m *mockPlatform) Architecture() string                                        { return "amd64" }
func (m *mockPlatform) Name() string                                                { return "macOS" }
func (m *mockPlatform) GetDataDir() string                                          { return "/tmp/data" }
func (m *mockPlatform) GetConfigDir() string                                        { return "/tmp/config" }
func (m *mockPlatform) GetCacheDir() string                                         { return "/tmp/cache" }
func (m *mockPlatform) GetLogDir() string                                           { return "/tmp/log" }
func (m *mockPlatform) GetIPCSocketPath() string                                    { return "/tmp/agentmgr.sock" }
func (m *mockPlatform) EnableAutoStart(ctx context.Context) error                   { return nil }
func (m *mockPlatform) DisableAutoStart(ctx context.Context) error                  { return nil }
func (m *mockPlatform) IsAutoStartEnabled(ctx context.Context) (bool, error)        { return false, nil }
func (m *mockPlatform) FindExecutable(name string) (string, error)                  { return "", nil }
func (m *mockPlatform) FindExecutables(name string) ([]string, error)               { return nil, nil }
func (m *mockPlatform) IsExecutableInPath(name string) bool                         { return m.executables[name] != "" }
func (m *mockPlatform) GetPathDirs() []string                                       { return nil }
func (m *mockPlatform) GetShell() string                                            { return "/bin/bash" }
func (m *mockPlatform) GetShellArg() string                                         { return "-c" }
func (m *mockPlatform) ShowNotification(title, message string) error                { return nil }
func (m *mockPlatform) ShowChangelogDialog(a, b, c, d string) platform.DialogResult { return 0 }

// ========== NPM Strategy Tests ==========

func TestNewNPMStrategy(t *testing.T) {
	plat := newMockPlatform()
	strategy := NewNPMStrategy(plat)

	if strategy == nil {
		t.Fatal("NewNPMStrategy returned nil")
	}
	if strategy.platform != plat {
		t.Error("platform not set correctly")
	}
}

func TestNPMStrategyName(t *testing.T) {
	strategy := NewNPMStrategy(newMockPlatform())
	if strategy.Name() != "npm" {
		t.Errorf("Name() = %q, want %q", strategy.Name(), "npm")
	}
}

func TestNPMStrategyMethod(t *testing.T) {
	strategy := NewNPMStrategy(newMockPlatform())
	if strategy.Method() != agent.MethodNPM {
		t.Errorf("Method() = %v, want %v", strategy.Method(), agent.MethodNPM)
	}
}

func TestNPMStrategyIsApplicable(t *testing.T) {
	t.Run("with npm available", func(t *testing.T) {
		plat := newMockPlatform()
		plat.executables["npm"] = "/usr/local/bin/npm"
		strategy := NewNPMStrategy(plat)

		if !strategy.IsApplicable(plat) {
			t.Error("IsApplicable should return true when npm is available")
		}
	})

	t.Run("without npm available", func(t *testing.T) {
		plat := newMockPlatform()
		strategy := NewNPMStrategy(plat)

		if strategy.IsApplicable(plat) {
			t.Error("IsApplicable should return false when npm is not available")
		}
	})
}

func TestExtractNPMPackageName(t *testing.T) {
	tests := []struct {
		command  string
		expected string
	}{
		{"npm install -g @anthropic-ai/claude-code", "@anthropic-ai/claude-code"},
		{"npm i -g claude-code", "claude-code"},
		{"npm install --global my-package", "my-package"},
		{"npm install -g package@1.2.3", "package"},
		{"npm install -g @scope/package@latest", "@scope/package"},
		{"npm -g install package", "install"}, // edge case: -g before install, takes first non-flag token
		{"npm install package", ""},           // no -g flag
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			result := extractNPMPackageName(tt.command)
			if result != tt.expected {
				t.Errorf("extractNPMPackageName(%q) = %q, want %q", tt.command, result, tt.expected)
			}
		})
	}
}

// ========== Binary Strategy Tests ==========

func TestNewBinaryStrategy(t *testing.T) {
	plat := newMockPlatform()
	strategy := NewBinaryStrategy(plat)

	if strategy == nil {
		t.Fatal("NewBinaryStrategy returned nil")
	}
	if strategy.platform != plat {
		t.Error("platform not set correctly")
	}
}

func TestBinaryStrategyName(t *testing.T) {
	strategy := NewBinaryStrategy(newMockPlatform())
	if strategy.Name() != "binary" {
		t.Errorf("Name() = %q, want %q", strategy.Name(), "binary")
	}
}

func TestBinaryStrategyMethod(t *testing.T) {
	strategy := NewBinaryStrategy(newMockPlatform())
	if strategy.Method() != agent.MethodNative {
		t.Errorf("Method() = %v, want %v", strategy.Method(), agent.MethodNative)
	}
}

func TestBinaryStrategyIsApplicable(t *testing.T) {
	plat := newMockPlatform()
	strategy := NewBinaryStrategy(plat)

	// Binary strategy should always be applicable
	if !strategy.IsApplicable(plat) {
		t.Error("IsApplicable should always return true for binary strategy")
	}
}

func TestExtractVersionFromOutput(t *testing.T) {
	tests := []struct {
		output   string
		expected string
	}{
		{"version 1.2.3", "1.2.3"},
		{"v1.2.3", "1.2.3"},
		{"1.2.3", "1.2.3"},
		{"claude-code version 1.0.5", "1.0.5"},
		{"aider v0.50.0", "0.50.0"},
		{"Version: 2.0.0-beta.1", "2.0.0-beta.1"},
		{"some text v3.4.5 more text", "3.4.5"},
		{"no version here", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.output, func(t *testing.T) {
			result := extractVersionFromOutput(tt.output)
			if result != tt.expected {
				t.Errorf("extractVersionFromOutput(%q) = %q, want %q", tt.output, result, tt.expected)
			}
		})
	}
}

// ========== Pip Strategy Tests ==========

func TestNewPipStrategy(t *testing.T) {
	plat := newMockPlatform()
	strategy := NewPipStrategy(plat)

	if strategy == nil {
		t.Fatal("NewPipStrategy returned nil")
	}
	if strategy.platform != plat {
		t.Error("platform not set correctly")
	}
}

func TestPipStrategyName(t *testing.T) {
	strategy := NewPipStrategy(newMockPlatform())
	if strategy.Name() != "pip" {
		t.Errorf("Name() = %q, want %q", strategy.Name(), "pip")
	}
}

func TestPipStrategyMethod(t *testing.T) {
	strategy := NewPipStrategy(newMockPlatform())
	if strategy.Method() != agent.MethodPip {
		t.Errorf("Method() = %v, want %v", strategy.Method(), agent.MethodPip)
	}
}

func TestPipStrategyIsApplicable(t *testing.T) {
	tests := []struct {
		name        string
		executables map[string]string
		expected    bool
	}{
		{"with pip", map[string]string{"pip": "/usr/bin/pip"}, true},
		{"with pip3", map[string]string{"pip3": "/usr/bin/pip3"}, true},
		{"with pipx", map[string]string{"pipx": "/usr/local/bin/pipx"}, true},
		{"with uv", map[string]string{"uv": "/usr/local/bin/uv"}, true},
		{"with all", map[string]string{"pip": "x", "pip3": "x", "pipx": "x", "uv": "x"}, true},
		{"with none", map[string]string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plat := newMockPlatform()
			plat.executables = tt.executables
			strategy := NewPipStrategy(plat)

			if strategy.IsApplicable(plat) != tt.expected {
				t.Errorf("IsApplicable() = %v, want %v", strategy.IsApplicable(plat), tt.expected)
			}
		})
	}
}

func TestExtractPipPackageName(t *testing.T) {
	tests := []struct {
		packageField string
		command      string
		expected     string
	}{
		{"aider-chat", "", "aider-chat"},
		{"", "pip install aider-chat", "aider-chat"},
		{"", "pip3 install package-name", "package-name"},
		{"", "pipx install aider-chat", "aider-chat"},
		{"", "uv tool install ruff", "ruff"},
		{"", "pip install package==1.2.3", "package"},
		{"", "pip install package>=1.0", "package"},
		{"", "pip install -U package", "package"},
		{"", "", ""},
	}

	for _, tt := range tests {
		name := tt.packageField
		if name == "" {
			name = tt.command
		}
		t.Run(name, func(t *testing.T) {
			result := extractPipPackageName(tt.packageField, tt.command)
			if result != tt.expected {
				t.Errorf("extractPipPackageName(%q, %q) = %q, want %q",
					tt.packageField, tt.command, result, tt.expected)
			}
		})
	}
}

func TestParseUVTextOutput(t *testing.T) {
	strategy := NewPipStrategy(newMockPlatform())

	tests := []struct {
		name     string
		output   string
		expected map[string]string // package name -> version
	}{
		{
			name:   "single package",
			output: "ruff v0.1.0",
			expected: map[string]string{
				"ruff": "0.1.0",
			},
		},
		{
			name:   "multiple packages",
			output: "ruff v0.1.0\naider-chat 0.50.0\nmypy v1.2.3",
			expected: map[string]string{
				"ruff":       "0.1.0",
				"aider-chat": "0.50.0",
				"mypy":       "1.2.3",
			},
		},
		{
			name:   "with empty lines",
			output: "\nruff v0.1.0\n\naider 0.50.0\n",
			expected: map[string]string{
				"ruff":  "0.1.0",
				"aider": "0.50.0",
			},
		},
		{
			name:   "with dashes (ignored)",
			output: "- ruff\n  v0.1.0\nreal-package v1.0.0",
			expected: map[string]string{
				"real-package": "1.0.0",
			},
		},
		{
			name:     "empty output",
			output:   "",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.parseUVTextOutput(tt.output)

			if len(result) != len(tt.expected) {
				t.Errorf("parseUVTextOutput returned %d packages, want %d", len(result), len(tt.expected))
			}

			for name, version := range tt.expected {
				pkg, found := result[name]
				if !found {
					t.Errorf("package %q not found in result", name)
					continue
				}
				if pkg.Version != version {
					t.Errorf("package %q version = %q, want %q", name, pkg.Version, version)
				}
			}
		})
	}
}

// ========== Brew Strategy Tests ==========

func TestNewBrewStrategy(t *testing.T) {
	plat := newMockPlatform()
	strategy := NewBrewStrategy(plat)

	if strategy == nil {
		t.Fatal("NewBrewStrategy returned nil")
	}
	if strategy.platform != plat {
		t.Error("platform not set correctly")
	}
}

func TestBrewStrategyName(t *testing.T) {
	strategy := NewBrewStrategy(newMockPlatform())
	if strategy.Name() != "brew" {
		t.Errorf("Name() = %q, want %q", strategy.Name(), "brew")
	}
}

func TestBrewStrategyMethod(t *testing.T) {
	strategy := NewBrewStrategy(newMockPlatform())
	if strategy.Method() != agent.MethodBrew {
		t.Errorf("Method() = %v, want %v", strategy.Method(), agent.MethodBrew)
	}
}

func TestBrewStrategyIsApplicable(t *testing.T) {
	tests := []struct {
		name        string
		platformID  platform.ID
		executables map[string]string
		expected    bool
	}{
		{"macOS with brew", platform.Darwin, map[string]string{"brew": "/opt/homebrew/bin/brew"}, true},
		{"macOS without brew", platform.Darwin, map[string]string{}, false},
		{"Linux with brew", platform.Linux, map[string]string{"brew": "/home/linuxbrew/.linuxbrew/bin/brew"}, true},
		{"Windows with brew", platform.Windows, map[string]string{"brew": "C:\\brew"}, false}, // brew not applicable on Windows
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plat := &mockPlatformWithID{
				mockPlatform: mockPlatform{executables: tt.executables},
				id:           tt.platformID,
			}
			strategy := NewBrewStrategy(plat)

			if strategy.IsApplicable(plat) != tt.expected {
				t.Errorf("IsApplicable() = %v, want %v", strategy.IsApplicable(plat), tt.expected)
			}
		})
	}
}

// mockPlatformWithID extends mockPlatform with configurable ID
type mockPlatformWithID struct {
	mockPlatform
	id platform.ID
}

func (m *mockPlatformWithID) ID() platform.ID { return m.id }

func TestExtractBrewPackageName(t *testing.T) {
	tests := []struct {
		packageField string
		command      string
		expected     string
	}{
		{"gh", "", "gh"},
		{"", "brew install gh", "gh"},
		{"", "brew install --cask visual-studio-code", "visual-studio-code"},
		{"", "brew install user/tap/formula", "formula"},
		{"", "brew install homebrew/core/package", "package"},
		{"", "brew install -q package", "package"},
		{"", "", ""},
	}

	for _, tt := range tests {
		name := tt.packageField
		if name == "" {
			name = tt.command
		}
		t.Run(name, func(t *testing.T) {
			result := extractBrewPackageName(tt.packageField, tt.command)
			if result != tt.expected {
				t.Errorf("extractBrewPackageName(%q, %q) = %q, want %q",
					tt.packageField, tt.command, result, tt.expected)
			}
		})
	}
}
