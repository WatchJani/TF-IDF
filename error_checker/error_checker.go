package error_checker

import "fmt"

func ErrorHandler(err error) {
	if err != nil{fmt.Println(err)}
}
