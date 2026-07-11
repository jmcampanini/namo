// Package namegen composes [prefix-]stamp-slug names.
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

// maxRounds bounds the unique-slug top-up loop so impossible counts fail
// with a clear error instead of spinning forever.
const maxRounds = 100

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
	// Count is the number of names to generate; must be at least 1.
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

// Generate returns Count names of the form [prefix-]stamp-slug, all sharing
// a single timestamp, with each slug unique within the batch.
func Generate(opts Options) ([]string, error) {
	if opts.Count < 1 {
		return nil, fmt.Errorf("count must be at least 1, got %d", opts.Count)
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

	names := make([]string, opts.Count)
	for i, slug := range unique {
		parts := make([]string, 0, 3)
		if opts.Prefix != "" {
			parts = append(parts, opts.Prefix)
		}
		if stamp != "" {
			parts = append(parts, stamp)
		}
		names[i] = strings.Join(append(parts, slug), "-")
	}
	return names, nil
}

// validSlug rejects malformed hotdiva2000 output: its embedded word lists
// each contain one empty entry, which occasionally yields slugs with a
// leading, trailing, or doubled hyphen.
func validSlug(s string) bool {
	return s != "" &&
		!strings.HasPrefix(s, "-") &&
		!strings.HasSuffix(s, "-") &&
		!strings.Contains(s, "--")
}

func hotdivaSlugs(n int, prefixThreshold, suffixThreshold float64) []string {
	return hotdiva2000.GenerateWithOptions(hotdiva2000.Options{
		Formatting:      hotdiva2000.FormatSlug,
		PrefixThreshold: prefixThreshold,
		Results:         n,
		SuffixThreshold: suffixThreshold,
	})
}
