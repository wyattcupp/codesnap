package ignore

import "regexp"

type IgnorePattern struct {
	Pattern *regexp.Regexp
	Negate  bool
	LineNo  int
	Line    string
}

type CodebaseIgnorePatterns struct {
	patterns []*IgnorePattern
}
