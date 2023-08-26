package tf_idf

import (
	"bufio"
	"errors"
	"log"
	"os"
	e "root/error_checker"
	"strconv"
	"strings"
)

// Character for remove from word
var Char = map[byte]bool{
	',':  true,
	'?':  true,
	'.':  true,
	'!':  true,
	':':  true,
	';':  true,
	'\'': true,
	'`':  true,
	')':  true,
	'(':  true,
	'<':  true,
	'>':  true,
	'=':  true,
	'/':  true,
	'"':  true,
	'-':  true,
}

var StopWords map[string]bool = make(map[string]bool, 665)

func Tokenization(text string) (map[string]float32, []string, int) {
	allWords, lenArray := allWords(text)

	tokenization, wordCounter, uniqWords := make(map[string]float32, lenArray/4), make(map[string]uint, lenArray/4), make([]string, 0, lenArray/4) // malo veci negko sto bi trebao biti capacity, ali trebalo bi biti zanemarivo!

	for _, word := range allWords {
		//parse word to readable format if string empty continue loop
		if wordChanger(&word) != nil {
			continue
		}
		// word = strings.TrimSpace(word)
		//check if is number or if is stop word
		if IsNumber(word) || StopWords[word] {
			continue
		}

		//if nit exist, add as uniq word :D
		if _, ok := wordCounter[word]; !ok {
			uniqWords = append(uniqWords, word)
		}

		//word count in document
		wordCounter[word]++

		//frequency of word in document
		tokenization[word] = float32(wordCounter[word]) / float32(lenArray)
	}

	return tokenization, uniqWords, len(tokenization)
}

// faster way
func allWords(text string) ([]string, uint) {
	var (
		counter uint
		start   int
	)

	for index := 0; index < len(text); index++ {
		if text[index] == ' ' {
			counter++
		}
	}

	words := make([]string, counter+1)

	counter = 0

	for index := 0; index < len(text); index++ {
		if text[index] == ' ' {
			words[counter] = text[start:index]
			start = index + 1
			counter++
		}
	}

	if start < len(text) {
		words[counter] = text[start:]
	}

	return words, counter
}

// make perfect word for tokenization, lower case word without (. , : ? ! ...)
func wordChanger(text *string) error {
	for index := 0; index < len(*text); index++ {
		if _, ok := Char[(*text)[index]]; ok {
			characterRemove(text, index)
		}
	}

	*text = strings.ToLower(*text)

	if len(*text) == 0 {
		return errors.New("Empty string")
	}
	// fmt.Println(*text)

	return nil
}

// remove , ? . ! : ; in words
func characterRemove(word *string, index int) {
	if len(*word)-1 < index {
		return
	}

	*word = ((*word)[:index]) + ((*word)[index+1:])

}

func IsNumber(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func StopWordsInit(path string) {
	file, err := os.Open(path)
	e.ErrorHandler(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Println(err)
			continue
		}

		StopWords[scanner.Text()] = true
	}
}
