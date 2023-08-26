package main

import (
	"fmt"
	"root/constants"
	t "root/tf_idf"

	db "root/database"
	e "root/error_checker"

	h "root/helper"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	r "root/routes"
)

func init() {
	t.StopWordsInit(constants.STOP_WORD_PATH)
}

func main() {
	// IDF := t.NewTF_IDF()

	// IDF.InitData()

	// myFileForRead := f.ReadFile("./blog/test_file")

	// start := time.Now()

	//+ data := IDF.TF_IDF(myFileForRead)

	// var keyValueList []KeyValue
	// for k, v := range data {
	// 	keyValueList = append(keyValueList, KeyValue{k, v})
	// }

	// sort.Slice(keyValueList, func(i, j int) bool {
	// 	return keyValueList[i].Value > keyValueList[j].Value
	// })

	// topKeys := make([]string, 0, 8)
	// for i := 0; i < 8 && i < len(keyValueList); i++ {
	// 	topKeys = append(topKeys, keyValueList[i].Key)
	// }

	// fmt.Println(time.Since(start))

	// fmt.Println("Best 5 tags:", topKeys)

	//====================================================================================

	store, err := db.Open(db.PostgresConfiguration())
	e.ErrorHandler(err)

	IDF := t.NewTF_IDF(store.DB)
	IDF.DatabaseInit()

	fmt.Println(IDF.GetAllWordDocument())

	app := router.New()
	r.TF_IDF(app, IDF)

	fasthttp.ListenAndServe(h.Port(), app.Handler)
}
