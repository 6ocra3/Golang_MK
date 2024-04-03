package words

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/jdkato/prose/v2"
	"github.com/kljensen/snowball"
)

// Функция для стемминга
func StemmString(input string) ([]string, error) {
	// Удаление знаков препинания
	input = strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) && r != '\'' && r != '-' {
			return -1
		}
		return r
	}, input)

	// Кусок кода, который обрабатывает строку, выделяет слова и характеризует каждое слово
	doc, err := prose.NewDocument(input)
	if err != nil {
		return nil, err
	}

	appended := make(map[string]bool)
	var filtered []string
	for _, tok := range doc.Tokens() {
		if checkWord(tok) {

			// Обработка слова
			stemmed, err := snowball.Stem(tok.Text, "english", true)
			if err != nil {
				fmt.Println("can't stem word:", err)
				continue
			}

			// Добавление в итоговый ответ только уникальных слов
			if !appended[stemmed] {
				filtered = append(filtered, stemmed)
				appended[stemmed] = true
			}

		}
	}

	return filtered, nil
}

func checkWord(tok prose.Token) bool {
	// forbiddenTags - тэги, который будут только у местоимений, чатсиц, предлогов и союзов
	forbiddenTags := []string{"DT", "IN", "CC", "PRP", "PRP$"}
	return !slices.Contains(forbiddenTags, tok.Tag)
}
