package collector

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	tiktoken "github.com/pkoukk/tiktoken-go"
	"github.com/wyattcupp/codebase-tool/internal/ignore"
)

var DefaultPatterns = []string{
	// sensitive/Config files
	"**/.env", "**/.env.local", "**/.env.*", "**/.config", "**/.ini",
	"**/*.properties", "**/.cfg", "**/.conf", "**/*.yaml", "**/*.yml",
	"**/*.pem", "**/*.ppk", "**/*.key", "**/*.keystore", "**/*.pfx", "**/*.p12",
	"**/*.db", "**/*.sqlite", "**/*.sqlite3", "**/*.bak", "**/*.backup",
	"**/*.dump", "**/*.sql", "**/wp-config.php", "**/settings.py",
	"**/application.properties", "**/web.config", "**/*.kdbx", "**/*.csv",
	"**/*.xls", "**/*.xlsx", "**/*.pdf", "**/.npmrc", "**/.pypirc",
	"**/.htpasswd", "**/id_rsa", "**/id_dsa", "**/.netrc", "**/.git-credentials",

	// binary/compiled files
	"**/*.exe", "**/*.dll", "**/*.so", "**/*.dylib", "**/*.bin", "**/*.obj",
	"**/*.o", "**/*.a", "**/*.lib", "**/*.pyc", "**/*.pyo", "**/*.pyd",
	"**/*.class", "**/*.jar", "**/*.war", "**/*.ear",

	// build outputs
	"**/*.map", "**/*.min.js", "**/*.min.css", "**/*.bundle.js", "**/*.bundle.css",

	// media/Assets
	"**/*.jpg", "**/*.jpeg", "**/*.png", "**/*.gif", "**/*.bmp", "**/*.ico",
	"**/*.svg", "**/*.webp", "**/*.mp3", "**/*.wav", "**/*.mp4", "**/*.avi",
	"**/*.mov", "**/*.flv", "**/*.webm", "**/*.ttf", "**/*.eot", "**/*.woff",
	"**/*.woff2", "**/*.otf",

	// compressed/archive
	"**/*.zip", "**/*.tar", "**/*.gz", "**/*.rar", "**/*.7z", "**/*.bz2",
	"**/*.xz",

	// common directories (with trailing slash to ensure directories)
	"**/.git/", "**/.svn/", "**/.hg/", "**/.bzr/",
	"**/.idea/", "**/.vscode/", "**/.vs/", "**/.eclipse/",
	"**/node_modules/", "**/vendor/", "**/bower_components/",
	"**/venv/", "**/env/", "**/.env/", "**/virtualenv/",
	"**/__pycache__/", "**/.pytest_cache/", "**/.mypy_cache/",
	"**/packages/", "**/jspm_packages/",
	"**/dist/", "**/build/", "**/out/", "**/target/",
	"**/bin/", "**/obj/", "**/Debug/", "**/Release/",
	"**/.next/", "**/.nuxt/", "**/.output/", "**/.vuepress/dist/",
	"**/docs/", "**/doc/", "**/documentation/", "**/javadoc/",
	"**/coverage/", "**/.coverage/", "**/htmlcov/",
	"**/test-results/", "**/test-reports/",
	"**/public/assets/", "**/static/assets/",
	"**/logs/", "**/log/", "**/tmp/", "**/temp/",
	"**/data/", "**/fixtures/", "**/uploads/", "**/.gitignore", "**/.codebase_ignore",
	"**/.Trash/", "**/$RECYCLE.BIN/", "**/System Volume Information/",
	"**/go.mod", "**/go.sum", "**/Gopkg.lock", "**/Gopkg.toml",

	// Logs/Cache/Temp
	"**/.cache", "**/.tmp", "**/.temp", "**/*.log", "**/*.logs",
	"**/.DS_Store", "**/Thumbs.db",

	// Lock files
	"**/*.lock", "**/package-lock.json", "**/yarn.lock", "**/Gemfile.lock",
	"**/poetry.lock",

	// Common metadata files
	"**/LICENSE", "**/NOTICE", "**/AUTHORS", "**/CONTRIBUTORS", "**/CHANGELOG",
	"**/LICENSE.*", "**/README.*", "**/CONTRIBUTING.*",
}

func GetTokenCount(text, model string) (int, error) {
	encoding, err := tiktoken.EncodingForModel(model)
	if err != nil {
		return 0, fmt.Errorf("failed to get encoding for model %s: %v", model, err)
	}

	tokens := encoding.Encode(text, nil, nil)
	return len(tokens), nil
}

func CollectCodebase(targetDir string, extraIgnores []string) (string, int64, error) {
	// 1. try to load lines from .codebase_ignore
	codebaseIgnoreFile := filepath.Join(targetDir, ".codebase_ignore")
	userLines := []string{}
	if data, err := os.ReadFile(codebaseIgnoreFile); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := scanner.Text()
			userLines = append(userLines, line)
		}
	}

	// 2. merge into a single CodebaseIgnorePatterns struct
	ignorePatterns := ignore.MergeIgnoreLines(
		DefaultPatterns,
		userLines,
		extraIgnores,
	)

	var sb strings.Builder

	// 3. walk the directory
	err := filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relPath, _ := filepath.Rel(targetDir, path)
		relPath = filepath.ToSlash(relPath)

		// 4. check if path is ignored
		isIgnored := ignorePatterns.MatchesPath(relPath)
		if isIgnored {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 5. if directory, we continue
		if d.IsDir() {
			return nil
		}

		// 6. MIME check for non-text in case something sneaky got through
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		buffer := make([]byte, 512)
		n, err := f.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		contentType := http.DetectContentType(buffer[:n])
		f.Close()
		if !strings.HasPrefix(contentType, "text/") {
			return nil
		}

		// 7. read entire file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sb.WriteString(fmt.Sprintf("**%s**\n", relPath))
		sb.WriteString("```")
		sb.WriteString("\n")
		sb.Write(data)
		sb.WriteString("\n```\n\n")
		return nil
	})

	if err != nil {
		return "", 0, err
	}

	result := sb.String()

	// tokenizer to retrieve token count
	tokenCount, err := GetTokenCount(result, "gpt-4o")

	return sb.String(), int64(tokenCount), nil
}
