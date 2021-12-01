package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

// siblingRunes находим предыдущую и следующую руну.
func siblingRunes(runes []rune, idx int) (rune, rune) {
	var prev rune
	var next rune
	if idx > 0 {
		prev = runes[idx-1]
	}
	if idx+1 < len(runes) {
		next = runes[idx+1]
	}
	return prev, next
}

// checkEscapedRune проверяет, что текущая руна экранированна символом /.
func checkEscapedRune(runes []rune, idx int) bool {
	escaped := false
	for i := range runes {
		prev, _ := siblingRunes(runes, i)
		if string(prev) != `\` {
			escaped = false
		} else {
			escaped = !escaped
		}
		if i == idx {
			break
		}
	}
	return escaped
}

func Unpack(val string) (string, error) {
	var newString strings.Builder
	runes := []rune(val)

	for i, currentRune := range runes {
		prevRune, nextRune := siblingRunes(runes, i)
		currentChr := string(currentRune)

		currentRuneIsDigit := unicode.IsDigit(currentRune)
		nextRuneIsDigit := unicode.IsDigit(nextRune)

		// Если символ \ не экранирован, то пропускаем его
		if currentChr == `\` && !checkEscapedRune(runes, i) {
			continue
		}

		if string(prevRune) == `\` {
			if currentRuneIsDigit {
				// Если предыдущий символ экранирует цифру, то это не цифра
				currentRuneIsDigit = !checkEscapedRune(runes, i)
			} else if currentChr != `\` {
				// Экранировать можно только цифры и \
				return "", ErrInvalidString
			}
		}

		if currentRuneIsDigit && nextRuneIsDigit {
			// Две подряд цифры - ошибка
			return "", ErrInvalidString
		}
		if currentRuneIsDigit {
			if prevRune == 0 {
				// Если цифра первая - то ошибка
				return "", ErrInvalidString
			}
			// Цифры пропускаем
			continue
		}

		if nextRuneIsDigit {
			repeats, err := strconv.Atoi(string(nextRune))
			if err != nil {
				return "", ErrInvalidString
			}
			newString.WriteString(
				strings.Repeat(currentChr, repeats),
			)
		} else {
			newString.WriteString(currentChr)
		}
	}
	return newString.String(), nil
}
