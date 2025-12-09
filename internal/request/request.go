package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"

	"github.com/jonnny013/go_html_server/internal/headers"
)

type parserState string

const (
	StateInit    parserState = "initialized"
	StateDone    parserState = "done"
	StateError   parserState = "error"
	StateHeaders parserState = "headers"
	StateBody    parserState = "body"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       parserState
	Headers     *headers.Headers
	Body        []byte
}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
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

func (r *Request) parse(data []byte, fileIsFinished bool) (int, error) {
	read := 0
outer:
	for {
		switch r.state {
		case StateError:
			return 0, fmt.Errorf("something went wrong")
		case StateDone:
			break outer

		case StateInit:

			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				r.state = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n

			r.state = StateHeaders

		case StateHeaders:

			n, done, err := r.Headers.Parse(data[read:])

			if err != nil {
				r.state = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			read += n

			if done {
				r.state = StateBody
			}
		case StateBody:

			contentLength := r.Headers.Get("content-length")

			if contentLength == "" {
				r.state = StateDone
				return 0, nil
			}

			length, err := strconv.Atoi(contentLength)

			if err != nil {
				r.state = StateError
				return 0, err
			}

			if length == 0 {
				r.state = StateDone
				return read, nil
			}

			remaining := min(length-len(r.Body), len(data[read:]))

			r.Body = append(r.Body, data[read:read+remaining]...)

			read += remaining

			if fileIsFinished && len(r.Body) != length {
				return 0, fmt.Errorf("body is incorrect length")
			}

			if len(data[read:]) == 0 {
				break outer
			}

			if len(r.Body) == length {
				r.state = StateDone
			}
		default:
			panic("we made a mistake as programmers")
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
		fileIsFinished := false
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, err
			}
			fileIsFinished = true
		}

		bufLen += n
		readN, err := req.parse(buf[:bufLen], fileIsFinished)
		if err != nil {
			return nil, err
		}
		if fileIsFinished {
			req.state = StateDone
			break
		}
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return req, nil
}
