package env

import "strings"

// PrefixTransformer adds or strips a prefix from env var keys.
type PrefixTransformer struct {
	prefix string
}

// NewPrefixTransformer creates a PrefixTransformer with the given prefix.
func NewPrefixTransformer(prefix string) *PrefixTransformer {
	return &PrefixTransformer{prefix: prefix}
}

// AddPrefix returns a new map with the prefix prepended to every key.
// Keys that already have the prefix are left unchanged.
func (t *PrefixTransformer) AddPrefix(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if strings.HasPrefix(k, t.prefix) {
			out[k] = v
		} else {
			out[t.prefix+k] = v
		}
	}
	return out
}

// StripPrefix returns a new map with the prefix removed from every key.
// Keys that do not have the prefix are dropped.
func (t *PrefixTransformer) StripPrefix(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if strings.HasPrefix(k, t.prefix) {
			out[strings.TrimPrefix(k, t.prefix)] = v
		}
	}
	return out
}

// RenamePrefix replaces oldPrefix with newPrefix on all matching keys.
// Keys without oldPrefix are passed through unchanged.
func RenamePrefix(secrets map[string]string, oldPrefix, newPrefix string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if strings.HasPrefix(k, oldPrefix) {
			out[newPrefix+strings.TrimPrefix(k, oldPrefix)] = v
		} else {
			out[k] = v
		}
	}
	return out
}
