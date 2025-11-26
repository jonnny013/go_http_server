package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
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

var SEPARATOR = "\r\n"

func parseRequestLine(b string) (*RequestLine, string, error) {
	idx := strings.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, b, nil
	}

	startLine := b[:idx]
	restOfMessage := b[idx+len(SEPARATOR):]

	parts := strings.Split(startLine, " ")

	if len(parts) != 3 {
		return nil, restOfMessage, fmt.Errorf("wrong length: %s", startLine)
	}

	httpParts := strings.Split(parts[2], "/")

	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, restOfMessage, fmt.Errorf("incorrect HTTP method: %s", parts[2])
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}

	if !rl.ValidMethod() || !rl.ValidPath() {
		return nil, restOfMessage, fmt.Errorf("malformed request line: %s", startLine)
	}

	return rl, restOfMessage, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	var req Request

	initialReq, err := io.ReadAll(reader)

	if err != nil {
		return &req, err
	}

	sl, _, err := parseRequestLine(string(initialReq))

	if err != nil {
		return &req, err
	}

	req.RequestLine = *sl

	return &req, err
}
