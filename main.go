package main

import (
	"fmt"
	f "root/file"
	t "root/tokenization"
)

func init() {
	t.StopWordsInit("./stop_words")
}

func main() {
	//one of blog
	text := f.ReadAllFile("./test_file")

	fmt.Println(t.Tokenization(text))
}
