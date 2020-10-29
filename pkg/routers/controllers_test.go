package routers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Test struct{}

var t *Test

func TestRoutes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Route Suite")
}

var _ = Describe("Route", func() {

	Context("when sending a request", func() {
		It("is GET /users", func() {
			resp, _, err := t.RunRequest("GET", "/users", "", nil)
			Expect(resp).To(Equal(200))
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("is POST /login", func() {
			// resp, content, err := t.RunRequest("POST", "/login", "", nil )
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
	r, err := http.NewRequest(method, "localhost:1234", p)

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
