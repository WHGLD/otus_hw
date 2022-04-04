package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")
var ErrInvalidIntParsing = errors.New("invalid integer parsing")

func Unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	var output strings.Builder

	var cur rune
	var prev rune

	var curStr string

	runeList := []rune(input)

	for i := 0; i < len(runeList); i++ {
		cur = runeList[i]
		if i > 0 {
			prev = runeList[i-1]
		}

		if unicode.IsDigit(cur) {
			if prev == 0 || unicode.IsDigit(prev) {
				return "", ErrInvalidString
			}

			runeIndex, e := strconv.Atoi(string(cur))
			if e != nil {
				return "", ErrInvalidIntParsing
			}

			if runeIndex == 0 {
				curStr = removeLastSymbol(output.String())
				output.Reset()
			} else {
				curStr = strings.Repeat(string(prev), runeIndex-1)
			}
		} else {
			curStr = string(cur)
		}

		output.WriteString(curStr)
	}

	return output.String(), nil
}

func removeLastSymbol(s string) string {
	runes := []rune(s)
	if lenOfString := len(runes); lenOfString > 1 {
		return string(runes[0 : lenOfString-1])
	}
	return ""
}
