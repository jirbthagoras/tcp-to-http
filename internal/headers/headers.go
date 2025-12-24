package headers

import (
	"bytes"
	"fmt"
)

var crlf = []byte("\r\n")
var separator = []byte(":")

type Headers map[string]string

func parseHeader(fieldLine []byte) (string, string, error) {
	fmt.Printf("%s\n", fieldLine)
	parts := bytes.SplitN(fieldLine, separator, 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field line")
	}

	key := parts[0]
	value := bytes.TrimSpace(parts[1])
	if bytes.HasSuffix(key, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}

	return string(key), string(value), nil
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
			read += len(separator)
			break
		}

		key, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}

		read += idx + len(separator)
		h[key] = value
	}

	return read, done, nil
}

func NewHeaders() Headers {
	return make(Headers)
}
