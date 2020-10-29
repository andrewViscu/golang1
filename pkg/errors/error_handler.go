package errors

import (
	"fmt"
)

//ErrNotNil l
func ErrNotNil(err error) {
	if err != nil {
		fmt.Printf("This is called from error_handler.go, %v\n", err)
	}
}
//TODO