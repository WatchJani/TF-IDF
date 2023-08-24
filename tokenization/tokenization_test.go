package tokenization

import (
	f "root/file"
	"testing"
)

func BenchmarkCharacterRemove(b *testing.B) {
	text := "Cow"

	for i := 0; i < b.N; i++ {
		characterRemove(&text, 3)
	}
}

func BenchmarkWordChanger(b *testing.B) {
	text := "industry."

	for i := 0; i < b.N; i++ {
		wordChanger(&text)
	}
}

var myText = f.ReadAllFile("./test_file")

func BenchmarkTokenization(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Tokenization(myText)
	}
}

func BenchmarkTokenization2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Tokenization2(myText)
	}
}

func BenchmarkAllWords(b *testing.B) {
	for i := 0; i < b.N; i++ {
		allWords(myText)
	}
}

//SVE STO JE IPOD SE VIŠE NE KORISTI :D
//========================================================================================

func BenchmarkTestGolangFunc(b *testing.B) {
	text := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."

	for i := 0; i < b.N; i++ {
		testGolangFunc(text)
	}
}

func BenchmarkTestSplitTextIntoWords(b *testing.B) {
	text := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."

	for i := 0; i < b.N; i++ {
		splitTextIntoWords(text)
	}
}
