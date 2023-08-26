package routes

import (
	t "root/tf_idf"

	"github.com/fasthttp/router"
)

func TF_IDF(app *router.Router, tf_idf *t.IDF) {

	app.POST("/training", tf_idf.Training)
	app.POST("/", tf_idf.GenerateTag)
}
