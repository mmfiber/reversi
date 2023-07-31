package strconverter

import "errors"

func IntToChar(i int) (string, error) {
	if i < 1 || i > 26 {
		return "", errors.New("invalid number")
	}
	return string('a' - 1 + i), nil
}

func IntToCapitalizedChar(i int) (string, error) {
	if i < 1 || i > 26 {
		return "", errors.New("invalid number")
	}
	return string('A' - 1 + i), nil
}

func CharToRune(char string) rune {
	return []rune(char)[0]
}
