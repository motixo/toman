package format

import (
	"strconv"
	"strings"
)

func FormatWithCommas(n int) string {
	s := strconv.Itoa(n)
	nLen := len(s)
	if nLen <= 3 {
		return s
	}

	var parts []string
	for nLen > 3 {
		parts = append([]string{s[nLen-3 : nLen]}, parts...)
		nLen -= 3
	}
	parts = append([]string{s[:nLen]}, parts...)
	return strings.Join(parts, ",")
}
