package main

import (
	"fmt"
	"root/constants"
	f "root/file"
	t "root/tf_idf"
	"sort"
	"time"
)

func init() {
	t.StopWordsInit(constants.STOP_WORD_PATH)
}

type KeyValue struct {
	Key   string
	Value float32
}

func main() {
	IDF := t.NewTF_IDF()

	IDF.InitData()

	myFileForRead := f.ReadFile("./blog/test_file")

	start := time.Now()

	data := IDF.TF_IDF(myFileForRead)

	var keyValueList []KeyValue
	for k, v := range data {
		keyValueList = append(keyValueList, KeyValue{k, v})
	}

	sort.Slice(keyValueList, func(i, j int) bool {
		return keyValueList[i].Value > keyValueList[j].Value
	})

	topKeys := make([]string, 0, 8)
	for i := 0; i < 8 && i < len(keyValueList); i++ {
		topKeys = append(topKeys, keyValueList[i].Key)
	}

	fmt.Println(time.Since(start))

	fmt.Println("Best 5 tags:", topKeys)
}
