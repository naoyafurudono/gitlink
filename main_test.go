package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestConvertToGitHubURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "GitHub SSH URL",
			input:    "git@github.com:user/repo.git",
			expected: "https://github.com/user/repo",
			wantErr:  false,
		},
		{
			name:     "GitHub SSH URL without .git",
			input:    "git@github.com:user/repo",
			expected: "https://github.com/user/repo",
			wantErr:  false,
		},
		{
			name:     "GitHub HTTPS URL",
			input:    "https://github.com/user/repo.git",
			expected: "https://github.com/user/repo",
			wantErr:  false,
		},
		{
			name:     "GitHub HTTPS URL without .git",
			input:    "https://github.com/user/repo",
			expected: "https://github.com/user/repo",
			wantErr:  false,
		},
		{
			name:     "GitHub git:// URL",
			input:    "git://github.com/user/repo.git",
			expected: "https://github.com/user/repo",
			wantErr:  false,
		},
		{
			name:     "GitLab SSH URL",
			input:    "git@gitlab.com:user/repo.git",
			expected: "https://gitlab.com/user/repo",
			wantErr:  false,
		},
		{
			name:     "Bitbucket SSH URL",
			input:    "git@bitbucket.org:user/repo.git",
			expected: "https://bitbucket.org/user/repo",
			wantErr:  false,
		},
		{
			name:     "Custom Git host SSH URL",
			input:    "git@git.example.com:user/repo.git",
			expected: "https://git.example.com/user/repo",
			wantErr:  false,
		},
		{
			name:     "Custom Git host HTTPS URL",
			input:    "https://git.example.com/user/repo.git",
			expected: "https://git.example.com/user/repo",
			wantErr:  false,
		},
		{
			name:     "Pepabo Git SSH URL",
			input:    "git@git.pepabo.com:team/project.git",
			expected: "https://git.pepabo.com/team/project",
			wantErr:  false,
		},
		{
			name:     "Invalid SSH URL format",
			input:    "git@invalidformat",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Unsupported URL format",
			input:    "ftp://example.com/repo.git",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "SSH URL with port",
			input:    "git@github.com:22:user/repo.git",
			expected: "https://github.com/22:user/repo",
			wantErr:  false,
		},
		{
			name:     "HTTP URL",
			input:    "http://github.com/user/repo.git",
			expected: "http://github.com/user/repo",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertToGitHubURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToGitHubURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("convertToGitHubURL() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReadFileLines(t *testing.T) {
	// Create a temporary file for testing
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	
	content := "line1\nline2\nline3\n"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		filePath string
		expected []string
		wantErr  bool
	}{
		{
			name:     "Read existing file",
			filePath: tmpFile,
			expected: []string{"line1", "line2", "line3"},
			wantErr:  false,
		},
		{
			name:     "Read non-existent file",
			filePath: filepath.Join(tmpDir, "nonexistent.txt"),
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines, err := readFileLines(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("readFileLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(lines) != len(tt.expected) {
					t.Errorf("readFileLines() returned %d lines, want %d", len(lines), len(tt.expected))
					return
				}
				for i, line := range lines {
					if line != tt.expected[i] {
						t.Errorf("readFileLines() line %d = %v, want %v", i, line, tt.expected[i])
					}
				}
			}
		})
	}
}

func TestReadFileLinesEmptyFile(t *testing.T) {
	// Test with empty file
	tmpDir := t.TempDir()
	emptyFile := filepath.Join(tmpDir, "empty.txt")
	
	err := os.WriteFile(emptyFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create empty test file: %v", err)
	}

	lines, err := readFileLines(emptyFile)
	if err != nil {
		t.Errorf("readFileLines() unexpected error for empty file: %v", err)
	}
	if len(lines) != 0 {
		t.Errorf("readFileLines() returned %d lines for empty file, want 0", len(lines))
	}
}

func TestReadFileLinesWithoutNewline(t *testing.T) {
	// Test file without trailing newline
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "no_newline.txt")
	
	content := "line1\nline2\nline3"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	lines, err := readFileLines(tmpFile)
	if err != nil {
		t.Errorf("readFileLines() unexpected error: %v", err)
	}
	
	expected := []string{"line1", "line2", "line3"}
	if len(lines) != len(expected) {
		t.Errorf("readFileLines() returned %d lines, want %d", len(lines), len(expected))
		return
	}
	
	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("readFileLines() line %d = %v, want %v", i, line, expected[i])
		}
	}
}

// Test git command functions
// These tests will only run if we're in a git repository

func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

func TestGetCurrentBranch(t *testing.T) {
	if !isGitRepository() {
		t.Skip("Not in a git repository")
	}

	branch, err := getCurrentBranch()
	if err != nil {
		t.Errorf("getCurrentBranch() error = %v", err)
		return
	}

	// Branch name should not be empty
	if branch == "" {
		t.Error("getCurrentBranch() returned empty string")
	}

	// Branch name should not contain newlines
	if strings.Contains(branch, "\n") {
		t.Error("getCurrentBranch() returned branch name with newline")
	}
}

func TestGetRemoteURL(t *testing.T) {
	if !isGitRepository() {
		t.Skip("Not in a git repository")
	}

	// Check if remote.origin.url is configured
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	if err := cmd.Run(); err != nil {
		t.Skip("No remote.origin.url configured")
	}

	url, err := getRemoteURL()
	if err != nil {
		t.Errorf("getRemoteURL() error = %v", err)
		return
	}

	// URL should not be empty
	if url == "" {
		t.Error("getRemoteURL() returned empty string")
	}

	// URL should not contain newlines
	if strings.Contains(url, "\n") {
		t.Error("getRemoteURL() returned URL with newline")
	}
}

func TestGetRepoRoot(t *testing.T) {
	if !isGitRepository() {
		t.Skip("Not in a git repository")
	}

	root, err := getRepoRoot()
	if err != nil {
		t.Errorf("getRepoRoot() error = %v", err)
		return
	}

	// Root should not be empty
	if root == "" {
		t.Error("getRepoRoot() returned empty string")
	}

	// Root should be an absolute path
	if !filepath.IsAbs(root) {
		t.Error("getRepoRoot() returned non-absolute path")
	}

	// Root should exist
	if _, err := os.Stat(root); os.IsNotExist(err) {
		t.Error("getRepoRoot() returned non-existent path")
	}

	// Root should contain .git directory
	gitDir := filepath.Join(root, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Error("getRepoRoot() returned path without .git directory")
	}
}

// Integration test for git functions working together
func TestGitFunctionsIntegration(t *testing.T) {
	if !isGitRepository() {
		t.Skip("Not in a git repository")
	}

	// Test that all git functions work without error
	branch, err := getCurrentBranch()
	if err != nil {
		t.Fatalf("getCurrentBranch() failed: %v", err)
	}

	root, err := getRepoRoot()
	if err != nil {
		t.Fatalf("getRepoRoot() failed: %v", err)
	}

	// Only test getRemoteURL if origin is configured
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	if err := cmd.Run(); err == nil {
		url, err := getRemoteURL()
		if err != nil {
			t.Fatalf("getRemoteURL() failed: %v", err)
		}

		// Test convertToGitHubURL with actual remote URL
		convertedURL, err := convertToGitHubURL(url)
		if err != nil {
			t.Logf("convertToGitHubURL() failed with actual remote URL %s: %v", url, err)
		} else {
			// Converted URL should be HTTP(S)
			if !strings.HasPrefix(convertedURL, "http://") && !strings.HasPrefix(convertedURL, "https://") {
				t.Errorf("convertToGitHubURL() returned non-HTTP(S) URL: %s", convertedURL)
			}
		}
	}

	t.Logf("Git integration test passed - branch: %s, root: %s", branch, root)
}