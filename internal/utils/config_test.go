package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/michaeldvinci/syllabus/internal/models"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		wantErr     bool
		expected    *models.Config
	}{
		{
			name: "valid config with single audiobook",
			yamlContent: `---
audiobooks:
  - title: "Test Book"
    audible: "https://www.audible.com/series/test"
    amazon: "https://www.amazon.com/dp/test"
`,
			wantErr: false,
			expected: &models.Config{
				Audiobooks: []models.Entry{
					{
						Title:   "Test Book",
						Audible: "https://www.audible.com/series/test",
						Amazon:  "https://www.amazon.com/dp/test",
					},
				},
			},
		},
		{
			name: "valid config with multiple audiobooks",
			yamlContent: `---
audiobooks:
  - title: "Book One"
    audible: "https://www.audible.com/series/book-one"
    amazon: "https://www.amazon.com/dp/book1"
  - title: "Book Two"
    audible: "https://www.audible.com/series/book-two"
    amazon: "https://www.amazon.com/dp/book2"
`,
			wantErr: false,
			expected: &models.Config{
				Audiobooks: []models.Entry{
					{
						Title:   "Book One",
						Audible: "https://www.audible.com/series/book-one",
						Amazon:  "https://www.amazon.com/dp/book1",
					},
					{
						Title:   "Book Two",
						Audible: "https://www.audible.com/series/book-two",
						Amazon:  "https://www.amazon.com/dp/book2",
					},
				},
			},
		},
		{
			name: "valid config with all fields",
			yamlContent: `---
audiobooks:
  - title: "Complete Book"
    audible: "https://www.audible.com/series/complete"
    amazon: "https://www.amazon.com/dp/complete"
    aud_num: 5
    aud_next: "Next Audio Title"
    aud_last: "Last Audio Title"
    amzn_num: 4
    amzn_next: "Next Amazon Title"
    amzn_last: "Last Amazon Title"
`,
			wantErr: false,
			expected: &models.Config{
				Audiobooks: []models.Entry{
					{
						Title:    "Complete Book",
						Audible:  "https://www.audible.com/series/complete",
						Amazon:   "https://www.amazon.com/dp/complete",
						AudNum:   5,
						AudNext:  "Next Audio Title",
						AudLast:  "Last Audio Title",
						AmznNum:  4,
						AmznNext: "Next Amazon Title",
						AmznLast: "Last Amazon Title",
					},
				},
			},
		},
		{
			name: "empty audiobooks list",
			yamlContent: `---
audiobooks: []
`,
			wantErr: false,
			expected: &models.Config{
				Audiobooks: []models.Entry{},
			},
		},
		{
			name: "malformed yaml - invalid syntax",
			yamlContent: `---
audiobooks:
  - title: "Test
    invalid yaml
`,
			wantErr: true,
		},
		{
			name: "malformed yaml - wrong structure",
			yamlContent: `---
not_audiobooks:
  - title: "Test Book"
`,
			wantErr: false,
			expected: &models.Config{
				Audiobooks: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test_config.yaml")
			
			err := os.WriteFile(tmpFile, []byte(tt.yamlContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Test LoadConfig
			config, err := LoadConfig(tmpFile)
			
			if tt.wantErr && err == nil {
				t.Errorf("LoadConfig() expected error but got none")
				return
			}
			
			if !tt.wantErr && err != nil {
				t.Errorf("LoadConfig() unexpected error: %v", err)
				return
			}
			
			if tt.wantErr {
				return // Expected error, test passed
			}

			// Compare results
			if config == nil {
				t.Errorf("LoadConfig() returned nil config")
				return
			}

			if len(config.Audiobooks) != len(tt.expected.Audiobooks) {
				t.Errorf("LoadConfig() audiobooks count = %d, want %d", 
					len(config.Audiobooks), len(tt.expected.Audiobooks))
				return
			}

			for i, entry := range config.Audiobooks {
				expected := tt.expected.Audiobooks[i]
				if entry.Title != expected.Title {
					t.Errorf("LoadConfig() entry[%d].Title = %q, want %q", 
						i, entry.Title, expected.Title)
				}
				if entry.Audible != expected.Audible {
					t.Errorf("LoadConfig() entry[%d].Audible = %q, want %q", 
						i, entry.Audible, expected.Audible)
				}
				if entry.Amazon != expected.Amazon {
					t.Errorf("LoadConfig() entry[%d].Amazon = %q, want %q", 
						i, entry.Amazon, expected.Amazon)
				}
				if entry.AudNum != expected.AudNum {
					t.Errorf("LoadConfig() entry[%d].AudNum = %v, want %v", 
						i, entry.AudNum, expected.AudNum)
				}
				if entry.AudNext != expected.AudNext {
					t.Errorf("LoadConfig() entry[%d].AudNext = %q, want %q", 
						i, entry.AudNext, expected.AudNext)
				}
				if entry.AudLast != expected.AudLast {
					t.Errorf("LoadConfig() entry[%d].AudLast = %q, want %q", 
						i, entry.AudLast, expected.AudLast)
				}
				if entry.AmznNum != expected.AmznNum {
					t.Errorf("LoadConfig() entry[%d].AmznNum = %v, want %v", 
						i, entry.AmznNum, expected.AmznNum)
				}
				if entry.AmznNext != expected.AmznNext {
					t.Errorf("LoadConfig() entry[%d].AmznNext = %q, want %q", 
						i, entry.AmznNext, expected.AmznNext)
				}
				if entry.AmznLast != expected.AmznLast {
					t.Errorf("LoadConfig() entry[%d].AmznLast = %q, want %q", 
						i, entry.AmznLast, expected.AmznLast)
				}
			}
		})
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent_file.yaml")
	if err == nil {
		t.Errorf("LoadConfig() expected error for nonexistent file but got none")
	}
}

func TestLoadConfig_InvalidPath(t *testing.T) {
	// Test with directory instead of file
	tmpDir := t.TempDir()
	_, err := LoadConfig(tmpDir)
	if err == nil {
		t.Errorf("LoadConfig() expected error for directory path but got none")
	}
}

func TestLoadConfig_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "empty.yaml")
	
	err := os.WriteFile(tmpFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create empty test file: %v", err)
	}

	config, err := LoadConfig(tmpFile)
	if err != nil {
		t.Errorf("LoadConfig() unexpected error for empty file: %v", err)
		return
	}

	if config == nil {
		t.Errorf("LoadConfig() returned nil config for empty file")
		return
	}

	if config.Audiobooks != nil {
		t.Errorf("LoadConfig() expected nil audiobooks for empty file, got %v", config.Audiobooks)
	}
}

func TestLoadConfig_PermissionDenied(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "no_permission.yaml")
	
	err := os.WriteFile(tmpFile, []byte("audiobooks: []"), 0000)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = LoadConfig(tmpFile)
	if err == nil {
		t.Errorf("LoadConfig() expected permission error but got none")
	}
}