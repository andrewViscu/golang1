package errors_test

import (
	. "andrewViscu/golang1/pkg/errors"
	"errors"
	"fmt"
	"testing"
)

func TestErrNotNil(t *testing.T) {
	err := errors.New("math: square root of negative number")
	go func(err error) {
		ErrNotNil(err)
	}(err)
	if err != nil {
		fmt.Println("Works!")
	}
}
