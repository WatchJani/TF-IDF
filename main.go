package main

import (
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
	store, err := db.Open(db.PostgresConfiguration())
	e.ErrorHandler(err)

	IDF := t.NewTF_IDF(store.DB)
	IDF.DatabaseInit()

	app := router.New()
	r.TF_IDF(app, IDF)

	fasthttp.ListenAndServe(h.Port(), app.Handler)
}
