package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	// db "andrewViscu/golang1/pkg/db"
	mw "andrewViscu/golang1/pkg/middleware"
)

func TestAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Suite")

	// fmt.Println("Everything worked!")
}

var (
	id string
)

var _ = Describe("Auth", func() {
	Context("When creating token from a random ID", func() {
		id = strconv.Itoa(rand.Int())
		token, err := mw.CreateToken(id)
		Expect(err).ShouldNot(HaveOccurred())

		It("should be equal to ID", func() {
			claims, err := mw.GetAuthenticatedUser(token)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(id).To(Equal(claims["user_id"]))
		})
	})
	Context("When sending request to auth required url", func() {
		req, err := http.NewRequest("POST", "/update/ "+id, nil)
		It("should get token", func() {
			Expect(err).ShouldNot(HaveOccurred())
			mw.GetToken(req)
		})
	})
	Context("Mocking middleware", func() {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		handlerToTest := mw.Handle(nextHandler)
		req := httptest.NewRequest("GET", "/users/"+id, nil)
		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	})
})
