package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	var req Request

	initialReq, err := io.ReadAll(reader)

	if err != nil {
		return &req, err
	}

	split := strings.Fields(string(initialReq))

	if len(split) < 3 {
		return &req, errors.New("empty request")
	}

	if split[0] != "GET" {
		return &req, fmt.Errorf("unsupported format: %s", split[0])
	}

	if split[1][0] != '/' {
		return &req, fmt.Errorf("unsupported format: %s", split[1])
	}

	if split[2] != "HTTP/1.1" {
		return &req, fmt.Errorf("unsupported format: %s", split[2])
	}

	req.RequestLine.Method = split[0]
	req.RequestLine.RequestTarget = split[1]
	req.RequestLine.HttpVersion = "1.1"

	return &req, nil
}
