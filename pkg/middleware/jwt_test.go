package middleware_test

import (
	"fmt"
	"testing"

	"strconv"
	"math/rand"
	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
	db "andrewViscu/golang1/pkg/db"
	mw "andrewViscu/golang1/pkg/middleware"
)



func TestAuth(t *testing.T) {
	
	id := strconv.Itoa(rand.Int())
	token, err := mw.CreateToken(id)
	errnotnil(err)
	claims, err := mw.GetAuthenticatedUser(token)
	errnotnil(err)
	fmt.Printf("Expected ID: %v, got %v\n", id, claims["user_id"])
	// fmt.Println("Everything worked!")
	
}

var _ = BeforeSuite(func() {
	db.DBConnect()
})


	

func errnotnil(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}





