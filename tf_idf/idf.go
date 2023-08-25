package tf_idf

import (
	"math"
	"sync"
)

type IDF struct {
	numberOfDocument float64
	allPossibleWords map[string]float64
}

func NewTF_IDF() *IDF {
	return &IDF{
		allPossibleWords: make(map[string]float64),
	}
}

func (i *IDF) SetWords(word string) {
	i.allPossibleWords[word]++
}

func (i *IDF) IncreaseNumDocument() {
	i.numberOfDocument++
}

func (i IDF) GetWordDocument(word string) float64 {
	return i.allPossibleWords[word]
}

func (i IDF) GetAllWordDocument() map[string]float64 {
	return i.allPossibleWords
}

func (i *IDF) TF_IDF(text string) map[string]float32 {
	tokenization, uniqWords, lenArray := Tokenization(text)

	tags := make(map[string]float32, lenArray)

	i.IncreaseNumDocument()

	for _, value := range uniqWords {
		i.SetWords(value)

		tags[value] = float32(math.Log10(i.numberOfDocument/i.allPossibleWords[value])) * tokenization[value]
	}

	return tags
}

func (i *IDF) IDF(text string) {
	_, uniqWords, _ := Tokenization(text)

	i.IncreaseNumDocument()

	for _, value := range uniqWords {
		i.SetWords(value)
	}
}

func (i *IDF) IDFSync(text string, wg *sync.WaitGroup) {
	_, uniqWords, _ := Tokenization(text)

	i.IncreaseNumDocument()

	for _, value := range uniqWords {
		i.SetWords(value)
	}

	wg.Done()
}
