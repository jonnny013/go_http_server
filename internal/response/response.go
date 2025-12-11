package response

import (
	"fmt"
	"io"

	"github.com/jonnny013/go_html_server/internal/headers"
)

type StatusCode int

const (
	StatusOk          StatusCode = 200
	StatusBadRequest  StatusCode = 400
	StatusSystemError StatusCode = 500
)

var crlf = []byte("\r\n")

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	status := []byte("HTTP/1.1 ")
	status = fmt.Appendf(status, "%v ", statusCode)
	switch statusCode {
	case StatusOk:
		status = fmt.Appendf(status, "OK")
	case StatusBadRequest:
		status = fmt.Appendf(status, "Bad Request")
	default:
		status = fmt.Appendf(status, "Internal Server Error")
	}

	status = append(status, crlf...)

	_, err := w.Write(status)
	return err
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()

	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, headers *headers.Headers) error {
	headersMap := headers.GetAll()

	for key, val := range headersMap {
		buf := make([]byte, 0, len(key)+len(val)+4)
		buf = fmt.Appendf(buf, "%s: %s", key, val)
		buf = append(buf, crlf...)

		if _, err := w.Write(buf); err != nil {
			return err
		}
	}

	_, err := w.Write(crlf)
	return err
}
