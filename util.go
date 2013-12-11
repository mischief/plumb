package plumb

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// bufio.Scanner function to split data by words and quoted strings
func scanStrings(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !unicode.IsSpace(r) {
			break
		}
	}

	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Scan until space, marking end of word.
	inquote := false
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if r == '\'' {
			inquote = !inquote
			continue
		}
		if unicode.IsSpace(r) && !inquote {
			return i + width, data[start:i], nil
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return 0, nil, nil
}

// parse attributes, unquoting values as necessary
func ParseAttr(line string) (Attr, error) {
	pa := make(map[string]string)

	scanw := bufio.NewScanner(strings.NewReader(line))
	scanw.Split(scanStrings)

	for scanw.Scan() {
		tpstr := scanw.Text()
		spl := strings.SplitN(tpstr, "=", 2)

		if len(spl) != 2 {
			return pa, fmt.Errorf("invalid tuple %q", tpstr)
		}

		spl[1] = strings.TrimLeft(spl[1], `'`)
		spl[1] = strings.TrimRight(spl[1], `'`)

		pa[spl[0]] = spl[1]

	}

	return pa, nil
}

// quote attribute value
func quote(str string) string {
	if !strings.ContainsAny(str, " '=\t") {
		return str
	}

	var out bytes.Buffer

	out.WriteRune('\'')

	for _, r := range str {
		out.WriteRune(r)
		if r == '\'' {
			out.WriteRune(r)
		}
	}

	out.WriteRune('\'')

	return out.String()
}
