package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var result strings.Builder
	var repeatBuffer string
	var isBackslashed bool

	if _, err := strconv.Atoi(str); err == nil {
		return "", ErrInvalidString
	}

	for i, value := range str {
		if i == 0 && unicode.IsDigit(value) {
			return "", ErrInvalidString
		}

		if isBackslashed {
			if !(unicode.IsDigit(value) || string(value) == string('\\')) {
				return "", ErrInvalidString
			}
			repeatBuffer = string(value)
			isBackslashed = false
			continue
		}

		if value == '\\' {
			result.WriteString(repeatBuffer)
			repeatBuffer = ""
			isBackslashed = true
			continue
		}

		if unicode.IsLetter(value) {
			result.WriteString(repeatBuffer)
			repeatBuffer = string(value)
			continue
		}

		if unicode.IsDigit(value) {
			if repeatBuffer == "" {
				return "", ErrInvalidString
			}

			repeatCount, err := strconv.Atoi(string(value))
			if err != nil {
				return "", err
			}

			result.WriteString(strings.Repeat(repeatBuffer, repeatCount))
			repeatBuffer = ""
		}
	}
	if isBackslashed {
		return "", ErrInvalidString
	}
	result.WriteString(repeatBuffer)
	return result.String(), nil
}
