package request

import (
	"bytes"
	"fmt"
	"io"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       parserState
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

var SEPARATOR = []byte("\r\n")
var ERR_MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request line")
var ERR_INCOMPLETE_REQUEST_LINE = fmt.Errorf("incomplete request line")
var ERR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")

type parserState string

const (
	StateInit parserState = "initialized"
	StateDone parserState = "done"
)

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERR_MALFORMED_REQUEST_LINE
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ERR_MALFORMED_REQUEST_LINE
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}

	return rl, read, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		switch r.state {
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n

			r.state = StateDone
		case StateDone:
			break outer
		}
	}

	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, 1024)
	bufLen := 0

	for !request.done() {
		// Mungut bytes satu persatu sesuai dengan jumlah numBytesPerRead
		// Sekarang buf harusnya berisi bytes dari reader.Data
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}

		// n adalah berapa jumlah bytes yang di copy ke buf?
		// Sehingga jumlah buf sekarang ditambah dengan num of copied bytes
		bufLen += n

		// Ngeparse, aku masih ndak tahu dia return apa, katane sih consumed bytes
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		// Kopi
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}
