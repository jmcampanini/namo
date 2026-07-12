// Package namegen composes [prefix-][stamp-]slug names.
package namegen

import (
	"fmt"
	"strings"
	"time"

	"charm.land/hotdiva2000"
)

// DefaultStampLayout is the fully sortable timestamp layout used when no
// stamp flags are given.
const DefaultStampLayout = "%y%m%d%H%M%S"

const (
	// MinCount is the minimum number of names Generate accepts.
	MinCount = 1
	// MaxCount is the maximum number of names Generate accepts.
	MaxCount = 100
	// maxRounds bounds the unique-slug top-up loop so impossible counts fail
	// with a clear error instead of spinning forever.
	maxRounds = 100
)

// Size selects how long generated slugs are.
type Size string

const (
	// SizeShort produces modifier-noun slugs with no extra words.
	SizeShort Size = "short"
	// SizeStandard occasionally adds extra prefix and suffix words.
	SizeStandard Size = "standard"
	// SizeLong always adds extra prefix and suffix words.
	SizeLong Size = "long"
)

// NormalizePrefix converts a prefix to lowercase ASCII words separated by dashes.
func NormalizePrefix(s string) (string, error) {
	var normalized strings.Builder
	normalized.Grow(len(s))
	separator := false

	for i := 0; i < len(s); i++ {
		char := s[i]
		if char >= 'A' && char <= 'Z' {
			char += 'a' - 'A'
		}
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
			if separator && normalized.Len() > 0 {
				normalized.WriteByte('-')
			}
			normalized.WriteByte(char)
			separator = false
			continue
		}
		separator = true
	}

	if normalized.Len() == 0 {
		return "", fmt.Errorf("prefix must contain at least one ASCII letter or digit")
	}
	return normalized.String(), nil
}

// ValidateCount validates the number of names requested for one batch.
func ValidateCount(count int) error {
	if count < MinCount || count > MaxCount {
		return fmt.Errorf("count must be between %d and %d, got %d", MinCount, MaxCount, count)
	}
	return nil
}

// ParseSize validates a user-supplied size name.
func ParseSize(s string) (Size, error) {
	switch Size(s) {
	case SizeShort, SizeStandard, SizeLong:
		return Size(s), nil
	}
	return "", fmt.Errorf("invalid size %q (valid: short, standard, long)", s)
}

func (s Size) thresholds() (prefix, suffix float64) {
	switch s {
	case SizeShort:
		return 0, 0
	case SizeLong:
		return 1, 1
	default:
		return 0.2, 0.2 // hotdiva2000 library defaults
	}
}

// Options configures Generate.
type Options struct {
	// Count is the number of names to generate; must be between MinCount and MaxCount.
	Count int
	// Now supplies the timestamp moment; nil means time.Now. It is called
	// exactly once per Generate so every name in a batch shares one stamp.
	Now func() time.Time
	// Prefix is an optional descriptive prefix joined with a hyphen.
	Prefix string
	// Size selects slug length; the zero value means SizeStandard.
	Size Size
	// Slugs supplies random slugs; nil means hotdiva2000.
	Slugs func(n int, prefixThreshold, suffixThreshold float64) []string
	// Stamp is a strftime-style layout; empty omits the timestamp.
	Stamp string
}

// Generate returns Count names of the form [prefix-][stamp-]slug, all sharing
// the same optional timestamp, with each slug unique within the batch.
func Generate(opts Options) ([]string, error) {
	if err := ValidateCount(opts.Count); err != nil {
		return nil, err
	}
	now := opts.Now
	if now == nil {
		now = time.Now
	}
	slugs := opts.Slugs
	if slugs == nil {
		slugs = hotdivaSlugs
	}

	var stamp string
	if opts.Stamp != "" {
		var err error
		stamp, err = formatStamp(now(), opts.Stamp)
		if err != nil {
			return nil, err
		}
	}

	prefixThreshold, suffixThreshold := opts.Size.thresholds()
	seen := make(map[string]struct{}, opts.Count)
	unique := make([]string, 0, opts.Count)
	for round := 0; len(unique) < opts.Count; round++ {
		if round == maxRounds {
			return nil, fmt.Errorf("generated only %d of %d unique slugs; reduce --count or use a larger --size", len(unique), opts.Count)
		}
		for _, slug := range slugs(opts.Count-len(unique), prefixThreshold, suffixThreshold) {
			if len(unique) == opts.Count {
				break
			}
			if !validSlug(slug) {
				continue
			}
			if _, dup := seen[slug]; dup {
				continue
			}
			seen[slug] = struct{}{}
			unique = append(unique, slug)
		}
	}

	parts := make([]string, 0, 2)
	if opts.Prefix != "" {
		parts = append(parts, opts.Prefix)
	}
	if stamp != "" {
		parts = append(parts, stamp)
	}
	namePrefix := strings.Join(parts, "-")
	if namePrefix != "" {
		namePrefix += "-"
	}

	names := make([]string, opts.Count)
	for i, slug := range unique {
		names[i] = namePrefix + slug
	}
	return names, nil
}

func validSlug(s string) bool {
	if s == "" {
		return false
	}
	componentHasContent := false
	for i := 0; i < len(s); i++ {
		char := s[i]
		switch {
		case (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9'):
			componentHasContent = true
		case char == '-' && componentHasContent:
			componentHasContent = false
		default:
			return false
		}
	}
	return componentHasContent
}

func hotdivaSlugs(n int, prefixThreshold, suffixThreshold float64) []string {
	return hotdiva2000.GenerateWithOptions(hotdiva2000.Options{
		Formatting:      hotdiva2000.FormatSlug,
		PrefixThreshold: prefixThreshold,
		Results:         n,
		SuffixThreshold: suffixThreshold,
	})
}
