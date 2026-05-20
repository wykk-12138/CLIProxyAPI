package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigOptional_ClaudeHeaderDefaults(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	configYAML := []byte(`
claude-header-defaults:
  user-agent: "  claude-cli/2.1.70 (external, cli)  "
  package-version: "  0.80.0  "
  runtime-version: "  v24.5.0  "
  os: "  MacOS  "
  arch: "  arm64  "
  timeout: "  900  "
  stabilize-device-profile: false
`)
	if err := os.WriteFile(configPath, configYAML, 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := LoadConfigOptional(configPath, false)
	if err != nil {
		t.Fatalf("LoadConfigOptional() error = %v", err)
	}

	if got := cfg.ClaudeHeaderDefaults.UserAgent; got != "claude-cli/2.1.70 (external, cli)" {
		t.Fatalf("UserAgent = %q, want %q", got, "claude-cli/2.1.70 (external, cli)")
	}
	if got := cfg.ClaudeHeaderDefaults.PackageVersion; got != "0.80.0" {
		t.Fatalf("PackageVersion = %q, want %q", got, "0.80.0")
	}
	if got := cfg.ClaudeHeaderDefaults.RuntimeVersion; got != "v24.5.0" {
		t.Fatalf("RuntimeVersion = %q, want %q", got, "v24.5.0")
	}
	if got := cfg.ClaudeHeaderDefaults.OS; got != "MacOS" {
		t.Fatalf("OS = %q, want %q", got, "MacOS")
	}
	if got := cfg.ClaudeHeaderDefaults.Arch; got != "arm64" {
		t.Fatalf("Arch = %q, want %q", got, "arm64")
	}
	if got := cfg.ClaudeHeaderDefaults.Timeout; got != "900" {
		t.Fatalf("Timeout = %q, want %q", got, "900")
	}
	if cfg.ClaudeHeaderDefaults.StabilizeDeviceProfile == nil {
		t.Fatal("StabilizeDeviceProfile = nil, want non-nil")
	}
	if got := *cfg.ClaudeHeaderDefaults.StabilizeDeviceProfile; got {
		t.Fatalf("StabilizeDeviceProfile = %v, want false", got)
	}
}

func TestLoadConfigOptional_ClaudeDeviceIDPerOS(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	configYAML := []byte(`
claude-header-defaults:
  device-id:
    macos: "  mac-device  "
    windows: "  windows-device  "
    default: "  default-device  "
  account-uuid: "  account-uuid  "
`)
	if err := os.WriteFile(configPath, configYAML, 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := LoadConfigOptional(configPath, false)
	if err != nil {
		t.Fatalf("LoadConfigOptional() error = %v", err)
	}

	if got := cfg.ClaudeHeaderDefaults.DeviceID.ValueForOS("MacOS"); got != "mac-device" {
		t.Fatalf("MacOS device-id = %q, want mac-device", got)
	}
	if got := cfg.ClaudeHeaderDefaults.DeviceID.ValueForOS("Windows"); got != "windows-device" {
		t.Fatalf("Windows device-id = %q, want windows-device", got)
	}
	if got := cfg.ClaudeHeaderDefaults.DeviceID.ValueForOS("Linux"); got != "default-device" {
		t.Fatalf("fallback device-id = %q, want default-device", got)
	}
	if got := cfg.ClaudeHeaderDefaults.AccountUUID; got != "account-uuid" {
		t.Fatalf("AccountUUID = %q, want account-uuid", got)
	}
}

func TestLoadConfigOptional_ClaudeDeviceIDLegacyString(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	configYAML := []byte(`
claude-header-defaults:
  device-id: "  legacy-device  "
`)
	if err := os.WriteFile(configPath, configYAML, 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := LoadConfigOptional(configPath, false)
	if err != nil {
		t.Fatalf("LoadConfigOptional() error = %v", err)
	}

	if got := cfg.ClaudeHeaderDefaults.DeviceID.Value(); got != "legacy-device" {
		t.Fatalf("legacy device-id = %q, want legacy-device", got)
	}
}
