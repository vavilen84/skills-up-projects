package http_client

import (
	"net/http"
)

type HTTPGetter interface {
	Get(url string) (resp *http.Response, err error)
}

type MockHTTPGetter struct {
	Resp *http.Response
	Err  error
}

func (m *MockHTTPGetter) Get(url string) (*http.Response, error) {
	return m.Resp, m.Err
}

func FetchData(getter HTTPGetter, url string) (*http.Response, error) {
	return getter.Get(url)
}
