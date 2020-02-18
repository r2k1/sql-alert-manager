package alert

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func humanizeDuration(duration time.Duration) string {
	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		singularName string
		pluralName   string
		amount       int64
	}{
		{"day", "days", days},
		{"hour", "hours", hours},
		{"minute", "minutes", minutes},
		{"second", "seconds", seconds},
	}

	var parts []string

	for _, chunk := range chunks {
		switch chunk.amount {
		case 0:
			continue
		case 1:
			parts = append(parts, fmt.Sprintf("%d %s", chunk.amount, chunk.singularName))
		default:
			parts = append(parts, fmt.Sprintf("%d %s", chunk.amount, chunk.pluralName))
		}
	}

	return strings.Join(parts, " ")
}
