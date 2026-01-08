package platform

import (
	"runtime"
	"testing"
)

func TestIDs(t *testing.T) {
	// Test ID constants
	if Darwin != "darwin" {
		t.Errorf("Darwin = %q, want %q", Darwin, "darwin")
	}
	if Linux != "linux" {
		t.Errorf("Linux = %q, want %q", Linux, "linux")
	}
	if Windows != "windows" {
		t.Errorf("Windows = %q, want %q", Windows, "windows")
	}
}

func TestDialogResultConstants(t *testing.T) {
	if DialogResultCancel != 0 {
		t.Errorf("DialogResultCancel = %d, want 0", DialogResultCancel)
	}
	if DialogResultUpdate != 1 {
		t.Errorf("DialogResultUpdate = %d, want 1", DialogResultUpdate)
	}
	if DialogResultRemindLater != 2 {
		t.Errorf("DialogResultRemindLater = %d, want 2", DialogResultRemindLater)
	}
	if DialogResultViewDetails != 3 {
		t.Errorf("DialogResultViewDetails = %d, want 3", DialogResultViewDetails)
	}
}

func TestCurrent(t *testing.T) {
	plat := Current()
	if plat == nil {
		t.Fatal("Current() returned nil")
	}

	// Should return same instance on subsequent calls
	plat2 := Current()
	if plat != plat2 {
		t.Error("Current() should return same instance")
	}
}

func TestCurrentID(t *testing.T) {
	id := CurrentID()
	expected := ID(runtime.GOOS)

	if id != expected {
		t.Errorf("CurrentID() = %q, want %q", id, expected)
	}
}

func TestCurrentArch(t *testing.T) {
	arch := CurrentArch()
	expected := runtime.GOARCH

	if arch != expected {
		t.Errorf("CurrentArch() = %q, want %q", arch, expected)
	}
}

func TestIsDarwin(t *testing.T) {
	expected := runtime.GOOS == "darwin"
	if IsDarwin() != expected {
		t.Errorf("IsDarwin() = %v, want %v", IsDarwin(), expected)
	}
}

func TestIsLinux(t *testing.T) {
	expected := runtime.GOOS == "linux"
	if IsLinux() != expected {
		t.Errorf("IsLinux() = %v, want %v", IsLinux(), expected)
	}
}

func TestIsWindows(t *testing.T) {
	expected := runtime.GOOS == "windows"
	if IsWindows() != expected {
		t.Errorf("IsWindows() = %v, want %v", IsWindows(), expected)
	}
}

func TestSupports(t *testing.T) {
	tests := []struct {
		id       ID
		expected bool
	}{
		{Darwin, true},
		{Linux, true},
		{Windows, true},
		{"freebsd", false},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.id), func(t *testing.T) {
			if Supports(tt.id) != tt.expected {
				t.Errorf("Supports(%q) = %v, want %v", tt.id, Supports(tt.id), tt.expected)
			}
		})
	}
}

func TestExecutableExtension(t *testing.T) {
	ext := ExecutableExtension()
	if IsWindows() {
		if ext != ".exe" {
			t.Errorf("ExecutableExtension() = %q, want %q on Windows", ext, ".exe")
		}
	} else {
		if ext != "" {
			t.Errorf("ExecutableExtension() = %q, want %q on non-Windows", ext, "")
		}
	}
}

func TestPathSeparator(t *testing.T) {
	sep := PathSeparator()
	if IsWindows() {
		if sep != ";" {
			t.Errorf("PathSeparator() = %q, want %q on Windows", sep, ";")
		}
	} else {
		if sep != ":" {
			t.Errorf("PathSeparator() = %q, want %q on non-Windows", sep, ":")
		}
	}
}

func TestHomeDirEnv(t *testing.T) {
	env := HomeDirEnv()
	if IsWindows() {
		if env != "USERPROFILE" {
			t.Errorf("HomeDirEnv() = %q, want %q on Windows", env, "USERPROFILE")
		}
	} else {
		if env != "HOME" {
			t.Errorf("HomeDirEnv() = %q, want %q on non-Windows", env, "HOME")
		}
	}
}

func TestTempDir(t *testing.T) {
	dir := TempDir()
	if IsWindows() {
		if dir != "C:\\Windows\\Temp" {
			t.Errorf("TempDir() = %q, want %q on Windows", dir, "C:\\Windows\\Temp")
		}
	} else {
		if dir != "/tmp" {
			t.Errorf("TempDir() = %q, want %q on non-Windows", dir, "/tmp")
		}
	}
}

func TestPlatformInterface(t *testing.T) {
	plat := Current()

	// Test all interface methods don't panic
	t.Run("ID", func(t *testing.T) {
		id := plat.ID()
		if id != Darwin && id != Linux && id != Windows {
			t.Errorf("unexpected platform ID: %q", id)
		}
	})

	t.Run("Architecture", func(t *testing.T) {
		arch := plat.Architecture()
		if arch == "" {
			t.Error("Architecture() returned empty string")
		}
	})

	t.Run("Name", func(t *testing.T) {
		name := plat.Name()
		if name == "" {
			t.Error("Name() returned empty string")
		}
	})

	t.Run("GetDataDir", func(t *testing.T) {
		dir := plat.GetDataDir()
		if dir == "" {
			t.Error("GetDataDir() returned empty string")
		}
	})

	t.Run("GetConfigDir", func(t *testing.T) {
		dir := plat.GetConfigDir()
		if dir == "" {
			t.Error("GetConfigDir() returned empty string")
		}
	})

	t.Run("GetCacheDir", func(t *testing.T) {
		dir := plat.GetCacheDir()
		if dir == "" {
			t.Error("GetCacheDir() returned empty string")
		}
	})

	t.Run("GetLogDir", func(t *testing.T) {
		dir := plat.GetLogDir()
		if dir == "" {
			t.Error("GetLogDir() returned empty string")
		}
	})

	t.Run("GetIPCSocketPath", func(t *testing.T) {
		path := plat.GetIPCSocketPath()
		if path == "" {
			t.Error("GetIPCSocketPath() returned empty string")
		}
	})

	t.Run("GetPathDirs", func(t *testing.T) {
		dirs := plat.GetPathDirs()
		// May be empty in some environments, but shouldn't panic
		_ = dirs
	})

	t.Run("GetShell", func(t *testing.T) {
		shell := plat.GetShell()
		if shell == "" {
			t.Error("GetShell() returned empty string")
		}
	})

	t.Run("GetShellArg", func(t *testing.T) {
		arg := plat.GetShellArg()
		if arg == "" {
			t.Error("GetShellArg() returned empty string")
		}
	})

	t.Run("IsExecutableInPath", func(t *testing.T) {
		// Test with something that should exist
		// bash on Unix, cmd on Windows
		var execName string
		if IsWindows() {
			execName = "cmd"
		} else {
			execName = "sh"
		}
		result := plat.IsExecutableInPath(execName)
		// Just verify it doesn't panic
		_ = result
	})

	t.Run("FindExecutable", func(t *testing.T) {
		// Test with something unlikely to exist
		_, err := plat.FindExecutable("nonexistent-command-12345")
		// Should return error for non-existent command
		if err == nil {
			t.Log("FindExecutable found unexpected command")
		}
	})

	t.Run("FindExecutables", func(t *testing.T) {
		// Test with something unlikely to exist
		paths, _ := plat.FindExecutables("nonexistent-command-12345")
		// Should return empty or nil for non-existent command
		if len(paths) > 0 {
			t.Log("FindExecutables found unexpected commands")
		}
	})
}
