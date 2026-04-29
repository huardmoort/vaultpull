package env

import "strings"

// FilterMode controls which keys are included or excluded.
type FilterMode int

const (
	FilterInclude FilterMode = iota
	FilterExclude
)

// Filter holds a set of key patterns and a mode (include or exclude).
type Filter struct {
	Mode     FilterMode
	Patterns []string
}

// NewIncludeFilter returns a Filter that keeps only keys matching any pattern.
func NewIncludeFilter(patterns []string) Filter {
	return Filter{Mode: FilterInclude, Patterns: patterns}
}

// NewExcludeFilter returns a Filter that drops keys matching any pattern.
func NewExcludeFilter(patterns []string) Filter {
	return Filter{Mode: FilterExclude, Patterns: patterns}
}

// Apply returns a new map containing only the entries allowed by the filter.
// Pattern matching is case-insensitive prefix/suffix/contains glob: a leading
// or trailing '*' is treated as a wildcard; otherwise exact match is used.
func (f Filter) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		matched := f.matchesAny(k)
		switch f.Mode {
		case FilterInclude:
			if matched {
				out[k] = v
			}
		case FilterExclude:
			if !matched {
				out[k] = v
			}
		}
	}
	return out
}

// matchesAny reports whether key matches at least one pattern.
func (f Filter) matchesAny(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range f.Patterns {
		if matchPattern(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

// matchPattern performs simple wildcard matching (* prefix/suffix/both).
func matchPattern(key, pattern string) bool {
	switch {
	case strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*"):
		return strings.Contains(key, strings.Trim(pattern, "*"))
	case strings.HasPrefix(pattern, "*"):
		return strings.HasSuffix(key, strings.TrimPrefix(pattern, "*"))
	case strings.HasSuffix(pattern, "*"):
		return strings.HasPrefix(key, strings.TrimSuffix(pattern, "*"))
	default:
		return key == pattern
	}
}
