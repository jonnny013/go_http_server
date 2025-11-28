package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) GetAll() map[string]string {
	return h.headers
}

func (h *Headers) Set(name, value string) {
	h.headers[strings.ToLower(name)] = value
}

var crlf = []byte("\r\n")

var errMalformedData = fmt.Errorf("data is malformed")

func isValidToken(str []byte) bool {
	for _, ch := range str {

		if unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) {
			continue
		} else {
			switch ch {
			case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
				continue
			default:
				return false
			}
		}

	}
	return true
}

func getKeyVal(data []byte) (string, string, error) {
	idx := bytes.IndexByte(data, ':')
	if idx == -1 {
		return "", "", errMalformedData
	}

	if data[idx-1] == ' ' {
		return "", "", errMalformedData
	}

	keyBytes := data[:idx]

	key := strings.TrimSpace(string(keyBytes))

	if key == "" || !isValidToken(keyBytes) {
		return "", "", errMalformedData
	}
	endIndex := bytes.Index(data, crlf)

	valueBytes := data[idx+1 : endIndex]

	value := strings.TrimSpace(string(valueBytes))

	if value == "" {
		return "", "", errMalformedData
	}

	return strings.ToLower(key), value, nil
}

func (h *Headers) Parse(data []byte) (n int, done bool, err error) {

	consumedBytes := 0
	done = false

	for {
		idx := bytes.Index(data[consumedBytes:], crlf)
		if idx == -1 {
			return 0, false, nil
		}
		if idx == 0 {
			done = true
			break
		}

		key, value, err := getKeyVal(data[consumedBytes:])

		if err != nil {
			return 0, false, err
		}

		h.Set(key, value)

		consumedBytes += idx + len(crlf)
	}
	return consumedBytes, done, nil
}
