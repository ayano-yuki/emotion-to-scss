package converter

import (
	"strings"

	"emotion-to-scss/internal/domain"
)

func Convert(style domain.Style) string {
	body := strings.Trim(style.CSS, "\r\n")
	var out strings.Builder
	out.WriteString(".")
	out.WriteString(style.ClassName)
	out.WriteString(" {\n")
	if strings.TrimSpace(body) != "" {
		for _, line := range strings.Split(body, "\n") {
			out.WriteString("  ")
			out.WriteString(strings.TrimRight(line, "\r"))
			out.WriteString("\n")
		}
	}
	out.WriteString("}\n")
	return out.String()
}

func ConvertAll(styles []domain.Style) string {
	var out strings.Builder
	for i, style := range styles {
		if i > 0 {
			out.WriteString("\n")
		}
		out.WriteString(Convert(style))
	}
	return out.String()
}
