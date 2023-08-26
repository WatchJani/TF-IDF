package tf_idf

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	e "root/error_checker"
	"sort"
	"sync"

	h "root/helper"

	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"
)

type KeyValue struct {
	Key   string
	Value float32
}

type IDF struct {
	*sqlx.DB
	numberOfDocument float64
	allPossibleWords map[string]float64
}

func NewTF_IDF(db *sqlx.DB) *IDF {
	return &IDF{
		db,
		0,
		make(map[string]float64),
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

//=================================================================================0

func (i *IDF) DatabaseInit() {
	query := "SELECT * FROM tag.dictionary;"
	rows, err := i.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			token     string
			frequency float64
		)

		if err := rows.Scan(&token, &frequency); err != nil {
			log.Fatal(err)
		}
		i.allPossibleWords[token] = frequency
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	query = "SELECT SUM(id) FROM tag.document;"
	var sum sql.NullFloat64
	err = i.QueryRow(query).Scan(&sum)
	if err != nil {
		log.Fatal(err)
	}

	i.numberOfDocument = sum.Float64
}

func (i *IDF) GenerateTag(ctx *fasthttp.RequestCtx) {
	text := ctx.PostBody()

	data := i.TF_IDF(string(text))

	var keyValueList []KeyValue
	for k, v := range data {
		keyValueList = append(keyValueList, KeyValue{k, v})
	}

	sort.Slice(keyValueList, func(i, j int) bool {
		return keyValueList[i].Value > keyValueList[j].Value
	})

	topKeys := make([]string, 0, 4)
	for i := 0; i < 4 && i < len(keyValueList); i++ {
		topKeys = append(topKeys, keyValueList[i].Key)
	}

	jsonData, err := json.Marshal(topKeys)
	if err != nil {
		log.Println("Greška prilikom konverzije u JSON:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetContentType("application/json")
	ctx.Response.SetBody(jsonData)
}

func (i *IDF) Training(ctx *fasthttp.RequestCtx) {
	text := ctx.PostBody()

	_, uniqWords, _ := Tokenization(string(text))

	var wg sync.WaitGroup
	wg.Add(3)

	go i.UpdateActiveState(uniqWords, &wg)
	go i.WriteToken(uniqWords, &wg)
	go i.WriteDocument(string(text), &wg)

	wg.Wait()

	ctx.SetContentType("application/json")
}

func (i *IDF) UpdateActiveState(uniqWords []string, wg *sync.WaitGroup) {
	i.IncreaseNumDocument()

	for _, value := range uniqWords {
		i.SetWords(value)
	}

	wg.Done()
}

func (i *IDF) WriteDocument(text string, wg *sync.WaitGroup) {
	_, err := i.Exec(h.InsertDocument(), text)
	e.ErrorHandler(err)

	wg.Done()
}

func (i *IDF) WriteToken(words []string, wg *sync.WaitGroup) {
	_, err := i.Exec(h.NewQuery(h.QueryGenerator(words)))
	e.ErrorHandler(err)

	wg.Done()
}
