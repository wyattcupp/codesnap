package collector

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// CollectCodebase recursively reads files from targetDir, ignoring anything
// that matches .codebase_ignore and any additional ignores. It returns a
// single markdown-formatted string that enumerates each file's contents.
func CollectCodebase(targetDir string, extraIgnores []string) (string, error) {
	ignoreRules, err := LoadIgnoreRules(filepath.Join(targetDir, ".codebase_ignore"))
	if err != nil {
		fmt.Printf("warn: failed to load .codebase_ignore (continuing without it): %v\n", err)
	}

	for _, ig := range extraIgnores {
		if ig != "" {
			ignoreRules = append(ignoreRules, ig)
		}
	}

	var sb strings.Builder
	var skipped []string

	err = filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		relPath, _ := filepath.Rel(targetDir, path)
		if strings.EqualFold(relPath, ".codebase_ignore") {
			return nil
		}

		if shouldIgnore(relPath, ignoreRules) {
			skipped = append(skipped, fmt.Sprintf("Ignored by rule or extension: %s", relPath))
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		buffer := make([]byte, 512)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		contentType := http.DetectContentType(buffer[:n])
		if !strings.HasPrefix(contentType, "text/") {
			skipped = append(skipped, fmt.Sprintf("Skipped non-text file: %s (MIME: %s)", relPath, contentType))
			return nil
		}

		file.Close()
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		unixRelPath := strings.ReplaceAll(relPath, `\`, `/`)

		sb.WriteString(fmt.Sprintf("**%s**\n", unixRelPath))
		sb.WriteString("```")
		sb.WriteString("\n")
		sb.Write(fileBytes)
		sb.WriteString("\n```\n\n")

		return nil
	})

	if err != nil {
		return "", err
	}

	if len(skipped) > 0 {
		fmt.Println("Skipped files:")
		for _, info := range skipped {
			fmt.Println(" -", info)
		}
	} else {
		fmt.Println("No files were skipped.")
	}

	return sb.String(), nil
}
