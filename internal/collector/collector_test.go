package collector

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTestFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file %s: %v", name, err)
	}
}

func createBinaryTestFile(t *testing.T, dir, name string, content []byte) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("Failed to create binary test file %s: %v", name, err)
	}
}

func TestCollectCodebase_IncludesTextFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "codebase-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// create txt files
	createTestFile(t, tempDir, "main.go", "package main\nfunc main() {}")
	createTestFile(t, tempDir, "README.md", "# Sample Project")

	// create a file that should be automatically skipped
	createBinaryTestFile(t, tempDir, "image.png", []byte{137, 80, 78, 71})

	output, err := CollectCodebase(tempDir, []string{})
	if err != nil {
		t.Fatalf("CollectCodebase failed: %v", err)
	}

	if !strings.Contains(output, "**main.go**") {
		t.Errorf("Expected main.go to be included in output")
	}
	if !strings.Contains(output, "**README.md**") {
		t.Errorf("Expected README.md to be included in output")
	}

	if strings.Contains(output, "image.png") {
		t.Errorf("Did not expect image.png to be included in output")
	}
}

func TestCollectCodebase_RespectsIgnoreRules(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "codebase-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	createTestFile(t, tempDir, "include.go", "package main")
	createTestFile(t, tempDir, "skip.go", "package skip")

	os.WriteFile(filepath.Join(tempDir, ".codebase_ignore"), []byte("skip.go"), 0644)

	output, err := CollectCodebase(tempDir, []string{})
	if err != nil {
		t.Fatalf("CollectCodebase failed: %v", err)
	}

	if !strings.Contains(output, "**include.go**") {
		t.Errorf("Expected include.go to be included in output")
	}
	if strings.Contains(output, "skip.go") {
		t.Errorf("Expected skip.go to be ignored")
	}
}

func TestCollectCodebase_SkipsNonTextFilesByMIME(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "codebase-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	binaryContent := []byte{0x00, 0x01, 0x02, 0x03}
	createBinaryTestFile(t, tempDir, "binaryfile", binaryContent)

	output, err := CollectCodebase(tempDir, []string{})
	if err != nil {
		t.Fatalf("CollectCodebase failed: %v", err)
	}

	if strings.Contains(output, "binaryfile") {
		t.Errorf("Expected binaryfile to be skipped based on MIME type")
	}
}

func TestShouldIgnore_ExcludesByExtension(t *testing.T) {
	ignoreRules := []string{}
	tests := []struct {
		relPath string
		expect  bool
	}{
		{"image.png", true},
		{"document.pdf", true},
		{"code.go", false},
	}
	for _, tt := range tests {
		got := shouldIgnore(tt.relPath, ignoreRules)
		if got != tt.expect {
			t.Errorf("shouldIgnore(%s) = %v; want %v", tt.relPath, got, tt.expect)
		}
	}
}

func TestMimeDetection_TextFile(t *testing.T) {
	content := []byte("This is a simple text file.")
	mimeType := http.DetectContentType(content)
	if !strings.HasPrefix(mimeType, "text/") {
		t.Errorf("Expected text/ MIME type, got %s", mimeType)
	}
}
