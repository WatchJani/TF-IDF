package helper

import (
	"fmt"
	"os"
)

func Port() string {
	return fmt.Sprintf(":%s", os.Getenv("PORT"))
}
