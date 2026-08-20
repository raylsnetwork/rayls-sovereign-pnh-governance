package database

import (
	"os"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeSSLMode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "should add sslmode=disable when missing",
			input:    "host=localhost port=5432 user=test",
			expected: "host=localhost port=5432 user=test sslmode=disable",
		},
		{
			name:     "should keep existing sslmode=disable",
			input:    "host=localhost sslmode=disable",
			expected: "host=localhost sslmode=disable",
		},
		{
			name:     "should keep existing sslmode=require",
			input:    "host=localhost sslmode=require",
			expected: "host=localhost sslmode=require",
		},
		{
			name:     "should remove sslmode=enable and add sslmode=disable",
			input:    "host=localhost sslmode=enable port=5432",
			expected: "host=localhost  port=5432 sslmode=disable",
		},
		{
			name:     "should handle empty string",
			input:    "",
			expected: " sslmode=disable",
		},
		{
			name:     "should handle multiple sslmode=enable",
			input:    "host=localhost sslmode=enable user=test sslmode=enable",
			expected: "host=localhost  user=test  sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeSSLMode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFindProjectMigrationsPath(t *testing.T) {
	tests := []struct {
		name          string
		setupDirs     []string // Directories to create (relative to tmpDir)
		workDir       string   // Directory to change to (relative to tmpDir)
		expectError   bool
		errorContains string
	}{
		{
			name:        "should find migrations folder in project root",
			setupDirs:   []string{"migrations"},
			workDir:     "",
			expectError: false,
		},
		{
			name:        "should find migrations folder in parent directory",
			setupDirs:   []string{"migrations", "cmd/api"},
			workDir:     "cmd/api",
			expectError: false,
		},
		{
			name:        "should find migrations folder multiple levels up",
			setupDirs:   []string{"migrations", "cmd/api/controllers/test"},
			workDir:     "cmd/api/controllers/test",
			expectError: false,
		},
		{
			name:          "should return error when migrations folder not found",
			setupDirs:     []string{"some/other/dir"},
			workDir:       "some/other/dir",
			expectError:   true,
			errorContains: "could not find migrations folder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save current directory
			originalWd, err := os.Getwd()
			require.NoError(t, err)
			defer func() { _ = os.Chdir(originalWd) }()

			// Create temporary directory structure
			tmpDir := t.TempDir()

			// Create all required directories
			for _, dir := range tt.setupDirs {
				fullPath := filepath.Join(tmpDir, dir)
				err = os.MkdirAll(fullPath, 0o750)
				require.NoError(t, err)
			}

			// Change to working directory
			workPath := tmpDir
			if tt.workDir != "" {
				workPath = filepath.Join(tmpDir, tt.workDir)
			}
			err = os.Chdir(workPath)
			require.NoError(t, err)

			// Test
			path, err := findProjectMigrationsPath()

			if tt.expectError {
				require.Error(t, err)
				assert.Empty(t, path)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				expectedPath := filepath.Join(tmpDir, "migrations")
				assert.Equal(t, expectedPath, path)
			}
		})
	}

	t.Run("should work from actual project directory", func(t *testing.T) {
		// This test assumes we're running from the project root or a subdirectory
		// Save current directory
		originalWd, err := os.Getwd()
		require.NoError(t, err)
		defer func() { _ = os.Chdir(originalWd) }()

		// Find the actual project migrations path
		path, err := findProjectMigrationsPath()

		// If this test is run from the actual project, it should find migrations
		// If run in isolation (e.g., CI), it might not find it - that's ok
		if err == nil {
			assert.NotEmpty(t, path)
			assert.DirExists(t, path)

			// Verify it's actually a migrations directory
			info, err := os.Stat(path)
			require.NoError(t, err)
			assert.True(t, info.IsDir())
			assert.Equal(t, "migrations", filepath.Base(path))
		}
	})
}

func TestRunMigrations_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name             string
		setupMigrations  bool // Whether to create migrations directory
		createDummyFile  bool // Whether to create a dummy migration file
		connectionString string
		expectError      bool
		errorDescription string
	}{
		{
			name:             "should return error when migrations directory not found",
			setupMigrations:  false,
			createDummyFile:  false,
			connectionString: "host=localhost port=5432 user=test password=test dbname=test sslmode=disable",
			expectError:      true,
			errorDescription: "Should return an error when migrations folder doesn't exist",
		},
		{
			name:             "should return error with invalid connection string",
			setupMigrations:  true,
			createDummyFile:  true,
			connectionString: "host=invalid-host-that-does-not-exist port=9999 user=fake password=fake dbname=fake sslmode=disable",
			expectError:      true,
			errorDescription: "Should error with invalid connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save current directory
			originalWd, err := os.Getwd()
			require.NoError(t, err)
			defer func() { _ = os.Chdir(originalWd) }()

			// Setup temporary directory
			tmpDir := t.TempDir()

			// Create migrations directory if needed
			if tt.setupMigrations {
				migrationsDir := filepath.Join(tmpDir, "migrations")
				err = os.MkdirAll(migrationsDir, 0o750)
				require.NoError(t, err)

				// Create dummy migration file if needed
				if tt.createDummyFile {
					migrationFile := filepath.Join(migrationsDir, "001_test.up.sql")
					err = os.WriteFile(migrationFile, []byte("CREATE TABLE test (id serial);"), 0o600)
					require.NoError(t, err)
				}
			}

			// Change to temp directory
			err = os.Chdir(tmpDir)
			require.NoError(t, err)

			// Test
			err = RunMigrations(tt.connectionString)

			if tt.expectError {
				require.Error(t, err, tt.errorDescription)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
