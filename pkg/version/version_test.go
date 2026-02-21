package version

import (
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	v := String()
	if v == "" {
		t.Error("version string should not be empty")
	}

	// Should have format X.Y.ZZZ
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		t.Errorf("version should have 3 parts, got %d: %s", len(parts), v)
	}
}

func TestFull(t *testing.T) {
	full := Full()

	if !strings.Contains(full, "H-lang") {
		t.Errorf("full version should contain 'H-lang': %s", full)
	}

	if !strings.Contains(full, "hlc") {
		t.Errorf("full version should contain 'hlc': %s", full)
	}

	if !strings.Contains(full, String()) {
		t.Errorf("full version should contain version string: %s", full)
	}
}

func TestVersion_Constants(t *testing.T) {
	if Major != 0 {
		t.Errorf("initial major version should be 0, got %d", Major)
	}

	if Minor != 0 {
		t.Errorf("initial minor version should be 0, got %d", Minor)
	}

	if Patch != 2 {
		t.Errorf("patch version should be 2, got %d", Patch)
	}
}
