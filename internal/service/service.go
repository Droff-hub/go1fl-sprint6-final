package service

import (
	"strings"
	"unicode"

	"github.com/Yandex-Practicum/golf1-sprint6-final/pkg/morse"
)

func Convert(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil
	}

	// Проверяем, состоит ли строка только из символов Морзе (. - и пробелы)
	isMorse := true
	for _, ch := range input {
		if ch != '.' && ch != '-' && ch != ' ' && ch != '\n' && ch != '\r' {
			isMorse = false
			break
		}
	}

	// Если есть буквы — это точно текст
	hasLetter := false
	for _, ch := range input {
		if unicode.IsLetter(ch) {
			hasLetter = true
			break
		}
	}

	if hasLetter {
		isMorse = false
	}

	var result string
	if isMorse {
		result = morse.ToText(input)
	} else {
		result = morse.ToMorse(input)
	}

	return result, nil
}