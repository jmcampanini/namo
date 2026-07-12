package namegen

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

func fixedNow() time.Time {
	return time.Date(2026, 7, 11, 15, 45, 1, 0, time.UTC)
}

// stubSlugs returns a slug source that hands out the given values in order,
// at most n per call, and nothing once exhausted.
func stubSlugs(values ...string) func(n int, prefixThreshold, suffixThreshold float64) []string {
	next := 0
	return func(n int, _, _ float64) []string {
		out := make([]string, 0, n)
		for ; n > 0 && next < len(values); n-- {
			out = append(out, values[next])
			next++
		}
		return out
	}
}

func TestNormalizePrefix(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "canonical unchanged", input: "release-2026", want: "release-2026"},
		{name: "uppercase lowered", input: "ReleaseCandidate42", want: "releasecandidate42"},
		{name: "digits retained", input: "123ABC789", want: "123abc789"},
		{name: "separator runs collapsed", input: "one \t--💥...two", want: "one-two"},
		{name: "edges trimmed", input: "--- alpha beta !!!", want: "alpha-beta"},
		{name: "unicode separates ASCII", input: "café東京Build", want: "caf-build"},
		{name: "accented uppercase is not converted", input: "Éclair", want: "clair"},
		{name: "slashes and underscores replaced", input: "team_name/feature", want: "team-name-feature"},
		{name: "empty", input: "", wantErr: true},
		{name: "hyphens only", input: "---", wantErr: true},
		{name: "whitespace only", input: " \t\n", wantErr: true},
		{name: "punctuation only", input: "!@#$%^&*()", wantErr: true},
		{name: "unicode only", input: "東京é💥", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizePrefix(tt.input)
			if tt.wantErr {
				if err == nil || !strings.Contains(err.Error(), "prefix must contain at least one ASCII letter or digit") {
					t.Fatalf("NormalizePrefix(%q) error = %v, want ASCII validation error", tt.input, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("NormalizePrefix(%q) error = %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("NormalizePrefix(%q) = %q, want %q", tt.input, got, tt.want)
			}
			canonical, err := NormalizePrefix(got)
			if err != nil || canonical != got {
				t.Fatalf("NormalizePrefix(%q) = %q, %v; want canonical input unchanged", got, canonical, err)
			}
		})
	}
}

func TestNormalizePrefixMalformedUTF8(t *testing.T) {
	got, err := NormalizePrefix(string([]byte{'A', 0xff, 'B', 0xc0, 'C'}))
	if err != nil {
		t.Fatalf("NormalizePrefix() error = %v", err)
	}
	if got != "a-b-c" {
		t.Fatalf("NormalizePrefix() = %q, want %q", got, "a-b-c")
	}

	_, err = NormalizePrefix(string([]byte{0xff, 0xfe}))
	if err == nil || err.Error() != "prefix must contain at least one ASCII letter or digit" {
		t.Fatalf("NormalizePrefix() error = %v, want ASCII validation error", err)
	}
}

func TestGenerateCompose(t *testing.T) {
	tests := []struct {
		name string
		opts Options
		want string
	}{
		{
			name: "default layout",
			opts: Options{Count: 1, Stamp: DefaultStampLayout},
			want: "260711154501-alpha-bravo",
		},
		{
			name: "prefix",
			opts: Options{Count: 1, Prefix: "debug-output", Stamp: DefaultStampLayout},
			want: "debug-output-260711154501-alpha-bravo",
		},
		{
			name: "no stamp",
			opts: Options{Count: 1},
			want: "alpha-bravo",
		},
		{
			name: "prefix without stamp",
			opts: Options{Count: 1, Prefix: "session"},
			want: "session-alpha-bravo",
		},
		{
			name: "short stamp layout",
			opts: Options{Count: 1, Stamp: "%H%M"},
			want: "1545-alpha-bravo",
		},
		{
			name: "dashed layout",
			opts: Options{Count: 1, Stamp: "%Y-%m-%d"},
			want: "2026-07-11-alpha-bravo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.opts.Now = fixedNow
			tt.opts.Slugs = stubSlugs("alpha-bravo")
			got, err := Generate(tt.opts)
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if len(got) != 1 || got[0] != tt.want {
				t.Fatalf("Generate() = %v, want [%q]", got, tt.want)
			}
		})
	}
}

func TestGenerateCountValidation(t *testing.T) {
	for _, count := range []int{0, -1} {
		_, err := Generate(Options{Count: count})
		if err == nil || !strings.Contains(err.Error(), "count must be at least 1") {
			t.Fatalf("Generate(Count: %d) error = %v, want count validation error", count, err)
		}
	}
}

func TestGenerateSharedStamp(t *testing.T) {
	calls := 0
	clock := func() time.Time {
		calls++
		return fixedNow().Add(time.Duration(calls) * time.Second)
	}
	got, err := Generate(Options{
		Count: 5,
		Now:   clock,
		Slugs: stubSlugs("a-a", "b-b", "c-c", "d-d", "e-e"),
		Stamp: DefaultStampLayout,
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if calls != 1 {
		t.Fatalf("clock called %d times, want 1", calls)
	}
	for _, name := range got {
		if !strings.HasPrefix(name, "260711154502-") {
			t.Fatalf("name %q does not share the single stamp", name)
		}
	}
}

func TestGenerateUniquenessTopUp(t *testing.T) {
	got, err := Generate(Options{
		Count: 3,
		Slugs: stubSlugs("dup-slug", "dup-slug", "dup-slug", "fresh-one", "fresh-two"),
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	want := []string{"dup-slug", "fresh-one", "fresh-two"}
	if len(got) != len(want) {
		t.Fatalf("Generate() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Generate() = %v, want %v", got, want)
		}
	}
}

func TestGenerateUniquenessBound(t *testing.T) {
	alwaysDup := func(int, float64, float64) []string { return []string{"same-slug"} }
	_, err := Generate(Options{Count: 2, Slugs: alwaysDup})
	if err == nil || !strings.Contains(err.Error(), "unique slugs") {
		t.Fatalf("Generate() error = %v, want unique-slug bound error", err)
	}
}

func TestGenerateFiltersMalformedSlugs(t *testing.T) {
	got, err := Generate(Options{
		Count: 2,
		Slugs: stubSlugs("-bad", "bad-", "a--b", "", "good-slug", "another-fine"),
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	want := []string{"good-slug", "another-fine"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Generate() = %v, want %v", got, want)
		}
	}
}

func TestGenerateIgnoresExcessSlugs(t *testing.T) {
	overSupply := func(int, float64, float64) []string {
		return []string{"one-a", "two-b", "three-c", "four-d", "five-e"}
	}
	got, err := Generate(Options{Count: 2, Slugs: overSupply})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("Generate() returned %d names, want 2", len(got))
	}
}

func TestSizeThresholds(t *testing.T) {
	tests := []struct {
		name       string
		size       Size
		wantPrefix float64
		wantSuffix float64
	}{
		{name: "short", size: SizeShort, wantPrefix: 0, wantSuffix: 0},
		{name: "standard", size: SizeStandard, wantPrefix: 0.2, wantSuffix: 0.2},
		{name: "zero value defaults to standard", size: "", wantPrefix: 0.2, wantSuffix: 0.2},
		{name: "long", size: SizeLong, wantPrefix: 1, wantSuffix: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPrefix, gotSuffix float64
			capture := func(_ int, prefixThreshold, suffixThreshold float64) []string {
				gotPrefix, gotSuffix = prefixThreshold, suffixThreshold
				return []string{"a-b"}
			}
			if _, err := Generate(Options{Count: 1, Size: tt.size, Slugs: capture}); err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if gotPrefix != tt.wantPrefix || gotSuffix != tt.wantSuffix {
				t.Fatalf("thresholds = (%v, %v), want (%v, %v)", gotPrefix, gotSuffix, tt.wantPrefix, tt.wantSuffix)
			}
		})
	}
}

func TestParseSize(t *testing.T) {
	for _, valid := range []string{"short", "standard", "long"} {
		if _, err := ParseSize(valid); err != nil {
			t.Fatalf("ParseSize(%q) error = %v", valid, err)
		}
	}
	_, err := ParseSize("bogus")
	if err == nil || !strings.Contains(err.Error(), "short, standard, long") {
		t.Fatalf("ParseSize(bogus) error = %v, want listing valid sizes", err)
	}
}

func TestGenerateRealSlugs(t *testing.T) {
	shape := regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)+$`)
	got, err := Generate(Options{Count: 200, Size: SizeShort})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	seen := make(map[string]struct{}, len(got))
	for _, slug := range got {
		if !shape.MatchString(slug) {
			t.Fatalf("slug %q does not match expected shape", slug)
		}
		if _, dup := seen[slug]; dup {
			t.Fatalf("duplicate slug %q in batch", slug)
		}
		seen[slug] = struct{}{}
	}
}
