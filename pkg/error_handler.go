package pkg

import (
	"fmt"
)

// func NewError(w http.ResponseWriter, err string) {
// 	w.WriteHeader(http.StatusInternalServerError)
// 	w.Write([]byte(`{ "error": "` + err + `", "response": 500 }`))
// 	return
// }
func ErrNotNil(err error) {
	if err != nil {
		fmt.Printf("This is called from error_handler.go, %v\n", err)
	}
}
//TODO