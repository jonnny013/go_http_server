package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode"
)

type parserState string

const (
	StateInit parserState = "initialized"
	StateDone parserState = "done"
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

func (r *RequestLine) ValidHTTP() bool {
	return r.HttpVersion == "1.1"
}

func (r *RequestLine) ValidPath() bool {
	return r.RequestTarget[0] == '/'
}

func (r *RequestLine) ValidMethod() bool {
	for _, r := range r.Method {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func (r *Request) isDone() bool {
	return r.state == StateDone
}

var SEPARATOR = []byte("\r\n")

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))

	if len(parts) != 3 {
		return nil, read, fmt.Errorf("wrong length: %s", startLine)
	}

	httpParts := bytes.Split(parts[2], []byte("/"))

	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, fmt.Errorf("incorrect HTTP method: %s", parts[2])
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}

	if !rl.ValidMethod() || !rl.ValidPath() {
		return nil, 0, fmt.Errorf("malformed request line: %s", startLine)
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

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := newRequest()
	req.state = StateInit

	buf := make([]byte, 1024)
	bufLen := 0
	for !req.isDone() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.state = StateDone
				break
			}
			return nil, err
		}

		bufLen += n
		readN, err := req.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return req, nil
}
