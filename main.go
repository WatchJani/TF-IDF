package main

import (
	"fmt"
	"root/constants"
	f "root/file"
	t "root/tf_idf"
)

func init() {
	t.StopWordsInit(constants.STOP_WORD_PATH)
}

func main() {
	IDF := t.NewTF_IDF()

	IDF.InitData()

	myFileForRead := f.ReadFile("./blog/test_file")

	fmt.Println(IDF.TF_IDF(myFileForRead))

}
