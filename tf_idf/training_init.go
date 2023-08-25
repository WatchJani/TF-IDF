package tf_idf

import (
	"root/constants"
	f "root/file"
	"strings"
	"sync"
)

func (i *IDF) InitData() {
	paths := f.ReadFile(constants.LOAD_ALL_TEST_PATH)

	//init all blog for analysis
	text := f.DocumentInit(strings.Fields(paths))

	var wg sync.WaitGroup

	//one of blog for training model
	for text := range text {
		wg.Add(1)

		go i.IDFSync(text, &wg)
	}

	wg.Wait()
}
