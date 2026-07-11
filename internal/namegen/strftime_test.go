package namegen

import (
	"strings"
	"testing"
	"time"
)

func TestFormatStamp(t *testing.T) {
	at := time.Date(2026, 7, 11, 15, 45, 1, 0, time.UTC)

	tests := []struct {
		name    string
		layout  string
		want    string
		wantErr string
	}{
		{name: "default layout", layout: "%y%m%d%H%M%S", want: "260711154501"},
		{name: "dashed date", layout: "%Y-%m-%d", want: "2026-07-11"},
		{name: "hour minute", layout: "%H%M", want: "1545"},
		{name: "seconds", layout: "%H%M%S", want: "154501"},
		{name: "escaped percent", layout: "%%", want: "%"},
		{name: "literal passthrough", layout: "v%y.%m", want: "v26.07"},
		{name: "literal digits stay literal", layout: "100%%", want: "100%"},
		{name: "empty layout", layout: "", want: ""},
		{name: "unsupported verb", layout: "%q", wantErr: "unsupported stamp verb %q"},
		{name: "trailing bare percent", layout: "%y%", wantErr: "ends with a bare %"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatStamp(at, tt.layout)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("formatStamp(%q) error = nil, want containing %q", tt.layout, tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("formatStamp(%q) error = %q, want containing %q", tt.layout, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("formatStamp(%q) error = %v", tt.layout, err)
			}
			if got != tt.want {
				t.Fatalf("formatStamp(%q) = %q, want %q", tt.layout, got, tt.want)
			}
		})
	}
}
