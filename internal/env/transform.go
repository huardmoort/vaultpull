package env

import (
	"strings"
)

// CaseTransform defines how key casing should be transformed.
type CaseTransform int

const (
	CaseNone  CaseTransform = iota
	CaseUpper               // UPPER_CASE
	CaseLower               // lower_case
)

// Transformer applies key and value transformations to a secrets map.
type Transformer struct {
	keyCase   CaseTransform
	valTrimWS bool
}

// NewTransformer returns a Transformer with the given options.
func NewTransformer(keyCase CaseTransform, trimValues bool) *Transformer {
	return &Transformer{
		keyCase:   keyCase,
		valTrimWS: trimValues,
	}
}

// Apply returns a new map with transformations applied to all keys and values.
// The original map is not mutated.
func (t *Transformer) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[t.transformKey(k)] = t.transformValue(v)
	}
	return out
}

func (t *Transformer) transformKey(k string) string {
	switch t.keyCase {
	case CaseUpper:
		return strings.ToUpper(k)
	case CaseLower:
		return strings.ToLower(k)
	default:
		return k
	}
}

func (t *Transformer) transformValue(v string) string {
	if t.valTrimWS {
		return strings.TrimSpace(v)
	}
	return v
}
