package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var (
		sb         strings.Builder
		prevLetter string
		isPrevNum  bool
	)

	runes := []rune(str)

	if len(runes) == 0 {
		return "", nil
	}

	if unicode.IsNumber(runes[0]) {
		return "", ErrInvalidString
	}

	for _, r := range runes {
		num, err := strconv.Atoi(string(r))

		if err == nil {
			if isPrevNum {
				return "", ErrInvalidString
			}

			sb.WriteString(strings.Repeat(prevLetter, num))
			isPrevNum = true
		} else {
			if !isPrevNum {
				sb.WriteString(prevLetter)
			}

			prevLetter = string(r)
			isPrevNum = false
		}
	}

	if !isPrevNum {
		sb.WriteString(prevLetter)
	}

	return sb.String(), nil
}
