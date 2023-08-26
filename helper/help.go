package helper

import (
	"fmt"
	"strings"
)

func QueryGenerator(token []string) []string {
	var valueStrings []string = make([]string, 0, len(token))
	for _, word := range token {
		valueStrings = append(valueStrings, fmt.Sprintf("('%s', 1)", word))
	}

	return valueStrings
}

func NewQuery(fields []string) string {
	return fmt.Sprintf(`
		INSERT INTO tag.dictionary (word, tokenization)
		VALUES %s
		ON CONFLICT (word) DO UPDATE
		SET tokenization = dictionary.tokenization + 1;
	`, strings.Join(fields, ","))
}

func InsertDocument() string {
	return "INSERT INTO tag.document (document) VALUES ($1)"
}
