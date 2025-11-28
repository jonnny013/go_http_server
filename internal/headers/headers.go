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

	valueBytes := data[idx+1 : endIndex]

	value := strings.TrimSpace(string(valueBytes))

	if value == "" {
		return "", "", errMalformedData
	}

	return key, value, nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

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

		h[key] = value

		consumedBytes += idx + len(crlf)
	}
	return consumedBytes, done, nil
}

func NewHeaders() Headers {
	return Headers{}
}
