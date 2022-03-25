package hw02unpackstring

import (
    "strings"
    "unicode"
    "strconv"
	"errors"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}

    var output strings.Builder

    var cur rune = 0
    var prev rune = 0

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

            if runeIndex := parseInt(string(cur)); runeIndex == 0 {
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
    lenOfString := len(runes)
    if lenOfString > 1 {
        return string(runes[0 : lenOfString-1])
    }
    return ""
}

func parseInt(s string) int {
    res, _ := strconv.Atoi(s)
    return res
}
