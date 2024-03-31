package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/jdkato/prose/v2"
	"github.com/kljensen/snowball"
)

func main() {

	sFlag := flag.String("s", "", "Флаг `-s` используется для ввода строки")
	flag.Parse()

	if *sFlag != "" {

		input := *sFlag

		// Кусок кода, который обрабатывает строку, выделяет слова и характеризует каждое слово
		doc, err := prose.NewDocument(input)
		if err != nil {
			fmt.Println("Can't preprocess input for prose:", err)
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
				if _, ok := appended[stemmed]; !ok {
					filtered = append(filtered, stemmed)
					appended[stemmed] = true
				}

			}
		}
		result := strings.Join(filtered, " ")
		fmt.Println(result)
	} else {
		fmt.Println(`На вход не подана строка. Используйте -s и "" `)
	}
}

func checkWord(tok prose.Token) bool {
	// forbidden_tags - тэги, который будут только у местоимений, чатсиц, предлогов и союзов
	forbidden_tags := []string{"DT", "IN", "CC", "PRP", "PRP$"}
	for _, token_code := range forbidden_tags {
		if tok.Tag == token_code {
			return false
		}
	}
	return true
}
