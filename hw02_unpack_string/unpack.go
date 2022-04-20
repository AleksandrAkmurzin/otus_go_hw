package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const escapeString = `\`

func Unpack(inputString string) (string, error) {
	var resultString strings.Builder

	runes := []rune(inputString)
	lastIndex := len(runes) - 1

	if lastIndex < 0 {
		return "", nil
	}

	if unicode.IsDigit(runes[0]) {
		return "", ErrInvalidString
	}

	isJustRepeated := false
	for i := 0; i <= lastIndex; i++ {
		currentRune := runes[i]
		if i == lastIndex {
			resultString.WriteRune(currentRune)
			break
		}

		if isJustRepeated && unicode.IsDigit(currentRune) {
			return "", ErrInvalidString
		}

		nextRune := runes[i+1]

		if string(currentRune) == escapeString {
			if !unicode.IsDigit(nextRune) && string(nextRune) != escapeString {
				return "", ErrInvalidString
			}

			currentRune = nextRune
			i++
			if i == lastIndex {
				resultString.WriteRune(currentRune)
				break
			}
			nextRune = runes[i+1]
		}

		isJustRepeated = writeOrRepeatRune(&resultString, currentRune, nextRune)
		if isJustRepeated {
			i++
		}
	}

	return resultString.String(), nil
}

// Return if rune was repeated.
func writeOrRepeatRune(builder *strings.Builder, current rune, next rune) bool {
	if unicode.IsDigit(next) {
		repeatCount, _ := strconv.Atoi(string(next))
		builder.WriteString(strings.Repeat(string(current), repeatCount))
		return true
	}

	builder.WriteRune(current)
	return false
}
