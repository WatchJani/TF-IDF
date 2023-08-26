package helper

import "testing"

func BenchmarkGenerator(b *testing.B) {
	myText := []string{" as asd a", " sad ", "asd", "kow", "asdo", "asd as ", " asd ", " as asd a"}

	for i := 0; i < b.N; i++ {
		QueryGenerator(myText)
	}
}
