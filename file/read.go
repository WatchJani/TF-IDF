package file

import (
	"io"
	"os"
	"root/constants"
	e "root/error_checker"
)

func ReadFile(path string) string {
	return string(ContentOfFIle(path))
}

func ContentOfFIle(path string) []byte {
	file, err := os.Open(path)
	e.ErrorHandler(err)

	bytes, err := io.ReadAll(file)
	e.ErrorHandler(err)

	return bytes
}

func ReadAllFileSync(path string, read chan string) {
	read <- string(ContentOfFIle(path))
}

func DocumentInit(paths []string) <-chan string {
	allDocument := make(chan string)

	go func() {
		for _, path := range paths {
			ReadAllFileSync(constants.PATH_REP+path, allDocument)
		}
		close(allDocument)
	}()

	return allDocument
}
