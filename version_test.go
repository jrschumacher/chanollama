package chanollama

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	version := GetVersion()
	
	if version == "" {
		t.Error("Version should not be empty")
	}
	
	// Check that version follows semver format (basic check)
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		t.Errorf("Version should follow semantic versioning format, got %s", version)
	}
	
	// Version constant should match the function
	if Version != GetVersion() {
		t.Errorf("Version constant (%s) should match GetVersion() (%s)", Version, GetVersion())
	}
}

func TestVersionFormat(t *testing.T) {
	version := Version
	
	// Check it starts with a digit
	if len(version) == 0 || version[0] < '0' || version[0] > '9' {
		t.Errorf("Version should start with a digit, got %s", version)
	}
	
	// Should not contain spaces
	if strings.Contains(version, " ") {
		t.Errorf("Version should not contain spaces, got %s", version)
	}
}