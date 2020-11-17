package middleware_test

import (
	// "net/http"
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


var _ = Describe("Auth", func() {
	var (
		id string
		// token string
	)
	Context("When creating token from a random ID", func() {
		id = strconv.Itoa(rand.Int())
		ts, err := mw.CreateToken(id)
		Expect(err).ShouldNot(HaveOccurred())

		It("should be equal to ID", func() {
			claims, err := mw.GetAuthenticatedUser(ts.AccessToken)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(id).To(Equal(claims["user_id"]))
		})
	})
	// Context("When sending request to auth required url", func() {
	// 	req, err := http.NewRequest("POST", "/update/ "+id, nil)
	// 	Expect(err).ShouldNot(HaveOccurred())
	// 	It("should get token", func() {
	// 		_, err := mw.GetToken(req)
	// 		Expect(err).ShouldNot(HaveOccurred())
			
	// 	})
	// })
})
