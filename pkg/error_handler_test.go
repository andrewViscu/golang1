package pkg_test

import (
	pkg "andrewViscu/golang1/pkg"
	"testing"
	"fmt"
	"errors"
)

func TestErrNotNil(t *testing.T)  {
	err := errors.New("math: square root of negative number")
	go func(err error) {
		pkg.ErrNotNil(err)
	}(err)
	if err != nil {
		fmt.Println("Works!")
	}
}

