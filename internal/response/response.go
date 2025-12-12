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
	StatusStream      StatusCode = 100
)

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer: writer,
	}
}

var crlf = []byte("\r\n")

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
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

	_, err := w.writer.Write(status)
	return err
}

func GetDefaultHeaders() *headers.Headers {
	h := headers.NewHeaders()

	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func GetContentLengthHeader(h *headers.Headers, contentLen int) {
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
}

func (w *Writer) WriteHeaders(headers *headers.Headers) error {
	headersMap := headers.GetAll()

	for key, val := range headersMap {
		buf := make([]byte, 0, len(key)+len(val)+4)
		buf = fmt.Appendf(buf, "%s: %s", key, val)
		buf = append(buf, crlf...)

		if _, err := w.writer.Write(buf); err != nil {
			return err
		}
	}

	_, err := w.writer.Write(crlf)
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBodyDone(dataLength int, hashSize []byte) (int, error) {

	i, err := w.WriteBody([]byte("0\r\n\r\n"))
	headers := headers.NewHeaders()
	headers.Set("X-Content-SHA256", fmt.Sprintf("%x", hashSize))
	headers.Set("X-Content-Length", fmt.Sprintf("%d", dataLength))
	w.WriteHeaders(headers)
	return i, err
}

func (w *Writer) WriteTrailers(h *headers.Headers) {
	h.Set("Trailer", "X-Content-SHA256 X-Content-Length")
}
