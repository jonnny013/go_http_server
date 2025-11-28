package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

var crlf = []byte("\r\n")

var errMalformedData = fmt.Errorf("data is malformed")

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

	if key == "" {
		return "", "", errMalformedData
	}
	endIndex := bytes.Index(data, crlf)

	valueBytes := data[idx+1:endIndex]

	value := strings.TrimSpace(string(valueBytes))

	if value == "" {
		return "", "", errMalformedData
	}

	return key, value, nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, crlf)

	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 0, true, nil
	}

	key, value, err := getKeyVal(data)

	if err != nil {
		return 0, false, err
	}

	h[key] = value

	return idx + len(crlf), false, nil

}

func NewHeaders() Headers {
	return Headers{}
}
