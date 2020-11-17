package routers_test

import (
	"andrewViscu/golang1/pkg/middleware"
	"andrewViscu/golang1/pkg/routers"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Test struct {
	Server *httptest.Server
}

type TestUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	At    string `json:"access_token"`
	Rt    string `json:"refresh_token"`
	City     string `json:"city"`
	Age      int    `json:"age"`
}

func TestRoutes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Route Suite")
}

func (test *Test) StartTestServer() {
	muxRouter := routers.ConfigureServer()
	test.Server = httptest.NewServer(muxRouter)
}

var _ = Describe("Route", func() {

	var (
		t  *Test
		tu TestUser
		id string
	)

	BeforeSuite(func() {
		t = &Test{}
		t.StartTestServer()
	})
	AfterSuite(func() {
		t.Server.Close()
	})
	Context("when sending a request", func() {
		It("is GET /", func() {
			resp, _, err := t.RunRequest("GET", "/", "", nil)
			ExpectDefault(resp, err)
		})
		It("is POST /register", func() {
			tu.Username = "Test2"
			tu.Password = "Test0000"
			resp, _, err := t.RunRequest("POST", "/register", "", tu)
			ExpectDefault(resp, err)
			log.Println("Created user")

		})
		It("is POST /login", func() {
			resp, content, err := t.RunRequest("POST", "/login", "", tu)
			ExpectDefault(resp, err)
			// log.Print(string(content))
			json.Unmarshal(content, &tu)
		})
		It("get user's ID", func() {
			// log.Print(tu)
			claims, err := middleware.GetAuthenticatedUser(tu.At)
			Expect(err).ShouldNot(HaveOccurred())
			id = fmt.Sprintf("%v", claims["user_id"])
		})
		It("is GET /users/:id", func() {
			resp, _, err := t.RunRequest("GET", "/users/"+id, tu.At, nil)
			ExpectDefault(resp, err)
			log.Println("Got user")
		})
		It("is PUT /users/:id", func() {
			var updateBody TestUser
			updateBody.Age = 13
			updateBody.City = "Kishinev"
			resp, _, err := t.RunRequest("PUT", "/users/"+id, tu.At, updateBody)
			ExpectDefault(resp, err)
			log.Println("Updated user")
		})
		It("is DELETE /users", func() {
			resp, _, err := t.RunRequest("DELETE", "/users/"+id, tu.At, nil)
			ExpectDefault(resp, err)
			log.Println("User deleted.")

		})
		It("is GET /users", func() {
			resp, _, err := t.RunRequest("GET", "/users", "", nil)
			ExpectDefault(resp, err)
			log.Println("Got all users")
		})
	})
})

func ExpectDefault(r *http.Response, err error) {
	Expect(err).ShouldNot(HaveOccurred())
	Expect(r.StatusCode).To(Equal(200))
}

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
				fmt.Println("DEBUGGING 5: ", err1)
				return nil, err1
			}
			p = strings.NewReader(string(pBytes))
		}
	}
	r, err := http.NewRequest(method, t.Server.URL+path, p)

	if err != nil {
		fmt.Println("DEBUGGING 4: ", err)
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")

	r.Header.Set("Authorization", "Bearer "+token)
	if err != nil {
		fmt.Println("DEBUGGING 3: ", err)
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
		fmt.Println("DEBUGGING 1: ", err)
		return
	}
	return RunRequest(req)
}
