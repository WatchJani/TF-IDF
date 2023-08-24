package tokenization

import (
	"bufio"
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
	// '<':  true,
	'>': true,
	'=': true,
	// '/':  true,
	'"': true,
	'-': true,
}

var StopWords map[string]bool = make(map[string]bool, 665)

func Tokenization(text string) map[string]uint {
	allWords := allWords(text)

	tokenization := make(map[string]uint, len(allWords))

	for _, word := range allWords {
		wordChanger(&word)

		if IsNumber(word) || StopWords[word] {
			continue
		}

		tokenization[word]++
	}

	return tokenization
}

// faster way
func allWords(text string) []string {
	var counter, start int

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

	return words
}

// make perfect word for tokenization, lower case word without (. , : ? ! ...)
func wordChanger(text *string) {
	for index := 0; index < len(*text); index++ {
		if _, ok := Char[(*text)[index]]; ok {
			characterRemove(text, index)
		}

		*text = strings.ToLower(*text)
	}
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

//SVE STO JE IPOD SE VIŠE NE KORISTI :D
//========================================================================================

// zasto je ovo 20x sporije??????
func Tokenization2(text string) map[string]uint {
	allWords := allWords(text)

	tokenization := make(map[string]uint, len(allWords))

	for _, word := range allWords {
		for index := 0; index < len(text); index++ {
			if _, ok := Char[(text)[index]]; ok {
				if len(word)-1 < index {
					continue
				}

				word = ((word)[:index]) + ((word)[index+1:])
			}

			text = strings.ToLower(text)
		}
		tokenization[word]++
	}

	return tokenization
}

// this is slower
func testGolangFunc(text string) []string {
	return strings.Fields(text)
}

// the most slower
func splitTextIntoWords(text string) []string {
	words := make([]string, 0)
	start := 0

	for index, char := range text {
		if char == ' ' {
			words = append(words, text[start:index])
			start = index + 1
		}
	}

	if start < len(text) {
		words = append(words, text[start:])
	}

	return words
}
