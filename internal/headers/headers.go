package headers

import (
	"bytes"
	"fmt"
	"strings"
)

var crlf = []byte("\r\n")
var separator = []byte(":")

type Headers struct {
	headers map[string]string
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Set(name string, value string) {
	name = strings.ToLower(name)
	if v, ok := h.headers[name]; ok {
		h.headers[name] = fmt.Sprintf("%s,%s", v, value)
	} else {
		h.headers[name] = value
	}

	h.headers[strings.ToLower(name)] = value
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func isToken(str []byte) bool {

	for _, ch := range str {
		found := false
		if ch >= 'a' && ch <= 'z' ||
			ch >= 'A' && ch <= 'Z' ||
			ch >= 0 && ch <= 9 {
			found = true
		}
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			found = true
		}

		if !found {
			return false
		}
	}

	return true
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, separator, 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field line")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])
	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}

	return string(name), string(value), nil
}

// n = byte consumed
// done = whether is done or not (the parsing I mean)
// err = the error of course
func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		// Cari registered nurse paling awal
		idx := bytes.Index(data[read:], crlf)
		if idx == -1 {
			break
		}

		// Header kosong wak
		if idx == 0 {
			done = true
			read += len(crlf)
			break
		}

		name, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}

		if !isToken([]byte(name)) {
			return 0, false, fmt.Errorf("malformed field name")
		}

		read += idx + len(crlf)
		h.Set(name, value)
	}

	return read, done, nil
}
