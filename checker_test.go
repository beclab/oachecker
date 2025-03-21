package oachecker

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCheckChart tests the CheckChart function
func TestCheckChart(t *testing.T) {
	// Test with valid chart path
	err := CheckChart("testdata/firefox")
	if err != nil {
		t.Errorf("CheckChart failed with valid chart path: %v", err)
	}
}

// TestCheckManifest tests the CheckManifest function
func TestCheckManifest(t *testing.T) {
	// Get a valid app configuration
	cfg, err := GetAppConfiguration("testdata/firefox")
	if err != nil {
		t.Fatalf("Failed to get app configuration: %v", err)
	}

	// Test with valid chart path and configuration
	err = CheckManifest("testdata/firefox", cfg)
	if err != nil {
		t.Errorf("CheckManifest failed with valid inputs: %v", err)
	}

	// Test with invalid configuration (modify a required field)
	invalidCfg := *cfg
	invalidCfg.ConfigVersion = ""
	err = CheckManifest("testdata/firefox", &invalidCfg)
	if err == nil {
		t.Error("CheckManifest should fail with invalid configuration")
	}
}

// TestLint tests the Lint function
func TestLint(t *testing.T) {
	// Test with valid chart path and default options
	err := Lint("testdata/firefox", DefaultLintOptions())
	if err != nil {
		t.Errorf("Lint failed with valid chart path and default options: %v", err)
	}

	// Test with custom options
	options := DefaultLintOptions().
		WithOwner("test-owner").
		WithAdmin("test-owner").
		WithCustomValidator(func(oacPath string, cfg *AppConfiguration) error {
			return nil // Always pass
		})
	err = Lint("testdata/firefox", options)
	if err != nil {
		t.Errorf("Lint failed with valid chart path and custom options: %v", err)
	}

	// Test with skipped manifest check
	options = DefaultLintOptions().SkipManifest()
	err = Lint("testdata/firefox", options)
	if err != nil {
		t.Errorf("Lint failed with skipped manifest check: %v", err)
	}

	// Test with skipped resource check
	options = DefaultLintOptions().SkipResources()
	err = Lint("testdata/firefox", options)
	if err != nil {
		t.Errorf("Lint failed with skipped resource check: %v", err)
	}

	// Test with invalid chart path
	err = Lint("testdata/nonexistent", DefaultLintOptions())
	if err == nil {
		t.Error("Lint should fail with invalid chart path")
	}
}

// TestLintWithDefaultOptions tests the LintWithDefaultOptions function
func TestLintWithDefaultOptions(t *testing.T) {
	// Test with valid chart path
	err := LintWithDefaultOptions("testdata/firefox")
	if err != nil {
		t.Errorf("LintWithDefaultOptions failed with valid chart path: %v", err)
	}

	// Test with invalid chart path
	err = LintWithDefaultOptions("testdata/nonexistent")
	if err == nil {
		t.Error("LintWithDefaultOptions should fail with invalid chart path")
	}
}

// TestLintWithSameOwnerAdmin tests the LintWithSameOwnerAdmin function
func TestLintWithSameOwnerAdmin(t *testing.T) {
	// Test with valid chart path and owner/admin
	err := LintWithSameOwnerAdmin("testdata/firefox", "test-owner-admin")
	if err != nil {
		t.Errorf("LintWithSameOwnerAdmin failed with valid inputs: %v", err)
	}

	// Test with invalid chart path
	err = LintWithSameOwnerAdmin("testdata/nonexistent", "test-owner-admin")
	if err == nil {
		t.Error("LintWithSameOwnerAdmin should fail with invalid chart path")
	}
}

// TestLintWithDifferentOwnerAdmin tests the LintWithDifferentOwnerAdmin function
func TestLintWithDifferentOwnerAdmin(t *testing.T) {
	// Test with valid chart path and different owner/admin
	err := LintWithDifferentOwnerAdmin("testdata/firefox", "test-owner", "test-admin")
	if err != nil {
		t.Errorf("LintWithDifferentOwnerAdmin failed with valid inputs: %v", err)
	}

	// Test with invalid chart path
	err = LintWithDifferentOwnerAdmin("testdata/nonexistent", "test-owner", "test-admin")
	if err == nil {
		t.Error("LintWithDifferentOwnerAdmin should fail with invalid chart path")
	}
}

// TestLintOptions tests the LintOptions methods
func TestLintOptions(t *testing.T) {
	// Test DefaultLintOptions
	options := DefaultLintOptions()
	if options.Owner != "" || options.Admin != "" || options.SkipManifestCheck || options.SkipResourceCheck || len(options.CustomValidators) != 0 {
		t.Error("DefaultLintOptions returned unexpected values")
	}

	// Test WithOwner
	options = DefaultLintOptions().WithOwner("test-owner")
	if options.Owner != "test-owner" {
		t.Errorf("WithOwner failed, expected 'test-owner', got '%s'", options.Owner)
	}

	// Test WithAdmin
	options = DefaultLintOptions().WithAdmin("test-admin")
	if options.Admin != "test-admin" {
		t.Errorf("WithAdmin failed, expected 'test-admin', got '%s'", options.Admin)
	}

	// Test WithSameOwnerAndAdmin
	options = DefaultLintOptions().WithSameOwnerAndAdmin("test-both")
	if options.Owner != "test-both" || options.Admin != "test-both" {
		t.Errorf("WithSameOwnerAndAdmin failed, expected both 'test-both', got Owner: '%s', Admin: '%s'", options.Owner, options.Admin)
	}

	// Test WithCustomValidator
	customValidator := func(oacPath string, cfg *AppConfiguration) error {
		return nil
	}
	options = DefaultLintOptions().WithCustomValidator(customValidator)
	if len(options.CustomValidators) != 1 {
		t.Errorf("WithCustomValidator failed, expected 1 validator, got %d", len(options.CustomValidators))
	}

	// Test WithAppDataValidator
	options = DefaultLintOptions()
	options.WithAppDataValidator()
	if len(options.CustomValidators) != 1 {
		t.Errorf("WithAppDataValidator failed, expected 1 validator, got %d", len(options.CustomValidators))
	}

	// Test SkipManifest
	options = DefaultLintOptions().SkipManifest()
	if !options.SkipManifestCheck {
		t.Error("SkipManifest failed, SkipManifestCheck should be true")
	}

	// Test SkipResources
	options = DefaultLintOptions().SkipResources()
	if !options.SkipResourceCheck {
		t.Error("SkipResources failed, SkipResourceCheck should be true")
	}
}

// Helper function to create a temporary test chart
func createTempTestChart(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "chart-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Copy the test chart to the temp directory
	err = copyDir("testdata/firefox", tempDir)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to copy test chart: %v", err)
	}

	return tempDir
}

// Helper function to copy a directory recursively
func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}
