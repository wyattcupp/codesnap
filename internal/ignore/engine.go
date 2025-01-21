package ignore

// some of this code was inspired by the go-gitignore project on GitHub //

import (
	"path/filepath"
	"regexp"
	"strings"
)

func CompileIgnoreLines(lines ...string) *CodebaseIgnorePatterns {
	patterns := &CodebaseIgnorePatterns{}
	for i, line := range lines {
		pattern, negate := getLineRegex(line)
		if pattern != nil {
			ip := &IgnorePattern{
				Pattern: pattern,
				Negate:  negate,
				LineNo:  i + 1,
				Line:    line,
			}
			patterns.patterns = append(patterns.patterns, ip)
		}
	}
	return patterns
}

func (gi *CodebaseIgnorePatterns) MatchesPath(f string) bool {
	match, _ := gi.GetMatchAndPattern(f)
	return match
}

func (gi *CodebaseIgnorePatterns) GetMatchAndPattern(f string) (bool, *IgnorePattern) {
	f = filepath.ToSlash(strings.TrimSpace(f))

	matchesPath := false
	var matchedPat *IgnorePattern
	for _, ignorePattern := range gi.patterns {
		if ignorePattern.Pattern.MatchString(f) {
			if !ignorePattern.Negate {
				matchesPath = true
				matchedPat = ignorePattern
			} else if matchesPath {
				matchesPath = false
				matchedPat = ignorePattern
			}
		}
	}
	return matchesPath, matchedPat
}

func getLineRegex(line string) (*regexp.Regexp, bool) {
	line = strings.TrimRight(line, "\r")

	// ignore comments
	if strings.HasPrefix(line, "#") {
		return nil, false
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return nil, false
	}

	negate := false
	if line[0] == '!' {
		negate = true
		line = line[1:]
	}

	if strings.HasPrefix(line, `\#`) || strings.HasPrefix(line, `\!`) {
		line = line[1:]
	}

	expr, _ := compileRegex(line)
	return expr, negate
}

func compileRegex(line string) (*regexp.Regexp, error) {
	// direcotries first
	if strings.HasPrefix(line, "/**/") {
		line = line[1:]
	}

	line = strings.ReplaceAll(line, "/**/", "(/|/.+/)")
	line = strings.ReplaceAll(line, "**/", "(|.*/)")
	line = strings.ReplaceAll(line, ".", `\.`)
	line = strings.ReplaceAll(line, "*", `[^/]*`)
	line = strings.ReplaceAll(line, "?", `[^/]`)

	// direcotyr suffix
	isDir := strings.HasSuffix(line, "/")
	if isDir {
		line = strings.TrimSuffix(line, "/")
		line += "/.*"
	}

	// anchors
	if strings.HasPrefix(line, "/") {
		line = "^" + line
	} else {
		line = "^(|.*?/)" + line
	}
	line += "$"

	return regexp.Compile(line)
}

func MergeIgnoreLines(sets ...[]string) *CodebaseIgnorePatterns {
	var all []string
	for _, set := range sets {
		all = append(all, set...)
	}
	return CompileIgnoreLines(all...)
}
