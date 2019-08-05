package mhttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestHTTPServer() *httptest.Server {
	response := []byte("Hello, ming-go/pkg!")

	return httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Write(response)
			},
		),
	)
}

func TestClientGetWithContext(t *testing.T) {
	ts := getTestHTTPServer()

	c := NewClient()
	resp, err := c.GetWithContext(context.Background(), ts.URL, nil, nil)
	assert.Nil(t, err)

	mhr, err := HttpResponseToMHttpResponse(resp, err)
	assert.Nil(t, err)

	assert.Equal(t, mhr.RespBody, []byte("Hello, ming-go/pkg!"))
}
