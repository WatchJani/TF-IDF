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

// constructor for TF-IDF
func NewTF_IDF(db *sqlx.DB) *IDF {
	return &IDF{
		db,
		0,
		make(map[string]float64),
	}
}

// set new word(increase) counter of word in one document
func (i *IDF) SetWords(word string) {
	i.allPossibleWords[word]++
}

// increase document
func (i *IDF) IncreaseNumDocument() {
	i.numberOfDocument++
}

// get number of document where is find word
func (i IDF) GetWordDocument(word string) float64 {
	return i.allPossibleWords[word]
}

// get all world with number in all document
func (i IDF) GetAllWordDocument() map[string]float64 {
	return i.allPossibleWords
}

// Return all analysis words [TF_IDF]
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

// init for database
func (i *IDF) DatabaseInit() {
	//load dictionary(all tokenization words)
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
	//load number of all document
	query = "SELECT SUM(id) FROM tag.document;"
	var sum sql.NullFloat64
	err = i.QueryRow(query).Scan(&sum)
	if err != nil {
		log.Fatal(err)
	}

	i.numberOfDocument = sum.Float64
}

// return to us 4 key word from text
func (i *IDF) GenerateTag(ctx *fasthttp.RequestCtx) {
	//get data from body (our text for analysis)
	text := ctx.PostBody()

	data := i.TF_IDF(string(text))

	//convert form map to struct for sorting
	var keyValueList []KeyValue
	for k, v := range data {
		keyValueList = append(keyValueList, KeyValue{k, v})
	}

	//sort our data
	sort.Slice(keyValueList, func(i, j int) bool {
		return keyValueList[i].Value > keyValueList[j].Value
	})

	//get top 4 word whose sorted
	topKeys := make([]string, 0, 4)
	for i := 0; i < 4 && i < len(keyValueList); i++ {
		topKeys = append(topKeys, keyValueList[i].Key)
	}

	//convert to json format
	jsonData, err := json.Marshal(topKeys)
	if err != nil {
		log.Println("Greška prilikom konverzije u JSON:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.Response.SetBody(jsonData)
}

// handler for training our model
func (i *IDF) Training(ctx *fasthttp.RequestCtx) {
	text := ctx.PostBody()

	//make tokenization
	_, uniqWords, _ := Tokenization(string(text))

	var wg sync.WaitGroup
	wg.Add(3)

	//update state
	go i.UpdateActiveState(uniqWords, &wg)
	//write tokenization word in database
	go i.WriteToken(uniqWords, &wg)
	//write our document in database
	go i.WriteDocument(string(text), &wg)

	wg.Wait()

	ctx.SetContentType("application/json")
}

// update our active state without database
func (i *IDF) UpdateActiveState(uniqWords []string, wg *sync.WaitGroup) {
	i.IncreaseNumDocument()

	for _, value := range uniqWords {
		i.SetWords(value)
	}

	wg.Done()
}

// wrote to data base our document for analysis
func (i *IDF) WriteDocument(text string, wg *sync.WaitGroup) {
	_, err := i.Exec(h.InsertDocument(), text)
	e.ErrorHandler(err)

	wg.Done()
}

// write to database our tokenization words
func (i *IDF) WriteToken(words []string, wg *sync.WaitGroup) {
	_, err := i.Exec(h.NewQuery(h.QueryGenerator(words)))
	e.ErrorHandler(err)

	wg.Done()
}
