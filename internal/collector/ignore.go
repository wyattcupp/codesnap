package collector

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

var excludedExtensions = []string{
	// Sensitive files
	".env", ".env.local", ".config", ".ini", ".properties", ".cfg", ".conf",
	".yaml", ".yml", ".pem", ".ppk", ".key", ".keystore", ".pfx", ".p12",
	".db", ".sqlite", ".bak", ".backup", ".dump", ".sql",
	"wp-config.php", "settings.py", "application.properties", "web.config",
	".kdbx", ".csv", ".xls", ".xlsx", ".pdf",
	".npmrc", ".pypirc", ".htpasswd", "id_rsa", "id_dsa", ".netrc", ".git-credentials",

	// Binary/Compiled files
	".exe", ".dll", ".so", ".dylib", ".bin", ".obj", ".o",

	// Image files
	".jpg", ".jpeg", ".png", ".gif", ".bmp", ".ico", ".svg", ".webp",

	// Audio/Video
	".mp3", ".wav", ".mp4", ".avi", ".mov", ".flv",

	// Compressed files
	".zip", ".tar", ".gz", ".rar", ".7z",

	// Cache/Build directories and files
	".cache", ".tmp", ".temp", ".log",
	"node_modules", "vendor", "dist", "build",
	".git", ".gitignore", ".svn",

	// IDE specific
	".idea", ".vscode", ".vs",
	".iml", ".project", ".classpath",

	// Package lock files (usually very large)
	"package-lock.json", "yarn.lock", "Gemfile.lock",

	// Large data files
	".parquet", ".avro", ".pb",

	// Documentation (usually available online)
	".pdf", ".doc", ".docx", ".ppt", ".pptx",
}

// LoadIgnoreRules loads lines from the specified .codebase_ignore file.
// Each non-empty, non-comment line is appended to a slice of ignore patterns.
func LoadIgnoreRules(ignoreFilePath string) ([]string, error) {
	file, err := os.Open(ignoreFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rules []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// skip comments/empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		rules = append(rules, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

// shouldIgnore checks if the given file or directory matches any of the ignore rules
func shouldIgnore(relPath string, ignoreRules []string) bool {
	for _, ext := range excludedExtensions {
		if strings.HasSuffix(strings.ToLower(relPath), ext) {
			return true
		}
	}

	for _, rule := range ignoreRules {
		if strings.Contains(relPath, rule) {
			return true
		}
		if filepath.Base(relPath) == rule {
			return true
		}
	}
	return false
}
