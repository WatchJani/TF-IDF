package tf_idf

import (
	f "root/file"
	"testing"
)

func BenchmarkWordChanger(b *testing.B) {
	text := "       industry.     "
	
	for i := 0; i < b.N; i++ {
		wordChanger(&text)
	}

}

var myText = f.ReadFile("./test_file")

func BenchmarkTokenization(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Tokenization(myText)
	}
}

func BenchmarkAllWords(b *testing.B) {
	for i := 0; i < b.N; i++ {
		allWords(myText)
	}
}
