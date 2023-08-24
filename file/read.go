package file

import (
	"io"
	"os"
	e "root/error_checker"
)

func ReadAllFile(path string) string {
	file, err := os.Open(path)
	e.ErrorHandler(err)

	bytes, err := io.ReadAll(file)
	e.ErrorHandler(err)
	return string(bytes)
}
