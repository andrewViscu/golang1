package routers_test

import (
	"andrewViscu/golang1/pkg/middleware"
	"andrewViscu/golang1/pkg/routers"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Test struct {
	Server *http.Server
}

type TestUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
	City     string `json:"city"`
	Age      int    `json:"age"`
}

var (
	t        *Test
	testUser TestUser
)

func TestRoutes(t *testing.T) {
	RegisterFailHandler(Fail)
	go routers.StartServer()
	RunSpecs(t, "Route Suite")
}

var _ = Describe("Route", func() {
	BeforeEach(func() {
		//Sorry, but I tried
		//routers.StartServer runs endlessly
		//without goroutines inside the router.go, I can't stop the server,
		//or it is possible but i'm tired of searching
	})
	AfterEach(func() {

	})
	Context("when sending a request", func() {
		It("is GET /users", func() {
			resp, _, err := t.RunRequest("GET", "/users", "", nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(200))
			// Expect(resp.StatusCode).To(Equal(200))
		})
		It("is POST /login", func() {
			testUser.Username = "Test"
			testUser.Password = "Test0000"
			resp, content, err := t.RunRequest("POST", "/login", "", testUser)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(200))

			json.Unmarshal(content, &testUser)
		})
		It("is PUT /users/:id", func() {
			var updateBody TestUser
			updateBody.Age = 13
			updateBody.City = "Kishinev"
			claims, err := middleware.GetAuthenticatedUser(testUser.Token)
			Expect(err).ShouldNot(HaveOccurred())
			id := fmt.Sprintf("%v", claims["user_id"])
			resp, content, err := t.RunRequest("PUT", "/users/"+id, testUser.Token, updateBody)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(200))
			fmt.Printf("Result: %v", string(content))
		})
	})
})

func (t *Test) GetRequest(method string, path string, token string, body interface{}) (*http.Request, error) {
	var p io.Reader

	if body != nil {
		switch body := body.(type) {
		case io.Reader:
			p = body
		case []byte:
			p = bytes.NewReader(body)
		case string:
			p = strings.NewReader(body)
		default:
			pBytes, err1 := json.Marshal(body)
			if err1 != nil {
				return nil, err1
			}
			p = strings.NewReader(string(pBytes))
		}
	}
	r, err := http.NewRequest(method, "http://localhost:1234"+path, p)

	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")

	r.Header.Set("Authorization", "Bearer "+token)
	if err != nil {
		return nil, err
	}
	return r, nil
}
func RunRequest(req *http.Request) (resp *http.Response, content []byte, err error) {

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err = client.Do(req)

	if err == nil {
		content, err = ioutil.ReadAll(resp.Body)

	}

	// todo right now, it doesn't hurt and makes things easier but it might cause issues in the future
	content = []byte(strings.TrimSpace(string(content)))

	return
}
func (t *Test) RunRequest(method string, path string, token string, body interface{}) (resp *http.Response, content []byte, err error) {
	var req *http.Request
	req, err = t.GetRequest(method, path, token, body)
	if err != nil {
		return
	}
	return RunRequest(req)
}
