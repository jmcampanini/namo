package namegen

import (
	"fmt"
	"strings"
	"time"
)

// formatStamp renders a strftime-style layout for t. Verbs are rendered
// individually rather than translated to a single Go layout string so
// literal characters can never collide with Go reference-time tokens.
func formatStamp(t time.Time, layout string) (string, error) {
	var b strings.Builder
	for i := 0; i < len(layout); i++ {
		if layout[i] != '%' {
			b.WriteByte(layout[i])
			continue
		}
		i++
		if i == len(layout) {
			return "", fmt.Errorf("stamp layout %q ends with a bare %%", layout)
		}
		switch layout[i] {
		case 'Y':
			b.WriteString(t.Format("2006"))
		case 'y':
			b.WriteString(t.Format("06"))
		case 'm':
			b.WriteString(t.Format("01"))
		case 'd':
			b.WriteString(t.Format("02"))
		case 'H':
			b.WriteString(t.Format("15"))
		case 'M':
			b.WriteString(t.Format("04"))
		case 'S':
			b.WriteString(t.Format("05"))
		case '%':
			b.WriteByte('%')
		default:
			return "", fmt.Errorf("unsupported stamp verb %%%c (supported: %%Y %%y %%m %%d %%H %%M %%S %%%%)", layout[i])
		}
	}
	return b.String(), nil
}
