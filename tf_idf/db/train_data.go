package db

import (
	h "root/helper"
	"sync"

	"github.com/jmoiron/sqlx"

	e "root/error_checker"
	t "root/tf_idf"
)

type Data struct {
	*sqlx.DB
}

func NewData(db *sqlx.DB) *Data {
	return &Data{
		db,
	}
}

func (d Data) TrainData(text string) {
	_, words, _ := t.Tokenization(text)

	var wg sync.WaitGroup
	wg.Add(2)

	go WriteToken(d, words, &wg)
	go WriteDocument(d, text, &wg)

	wg.Wait()
}

func WriteDocument(d Data, text string, wg *sync.WaitGroup) {
	_, err := d.Exec(h.InsertDocument(), text)
	e.ErrorHandler(err)

	wg.Done()
}

func WriteToken(d Data, words []string, wg *sync.WaitGroup) {
	_, err := d.Exec(h.NewQuery(h.QueryGenerator(words)))
	e.ErrorHandler(err)

	wg.Done()
}
