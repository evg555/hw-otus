package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var (
		sb         strings.Builder
		prevLetter string
		isPrevNum  bool
	)

	for i, letter := range str {
		num, err := strconv.Atoi(string(letter))

		if i == 0 && err == nil {
			return "", ErrInvalidString
		}

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

			prevLetter = string(letter)
			isPrevNum = false
		}
	}

	if !isPrevNum {
		sb.WriteString(prevLetter)
	}

	return sb.String(), nil
}
