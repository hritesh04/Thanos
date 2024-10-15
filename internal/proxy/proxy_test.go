package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hritesh04/thanos/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func init() {
	logger.InitLogger()
}

func TestProxyFunc(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	proxy := NewReverseProxy(server.URL)

	recorder := httptest.NewRecorder()

	request := httptest.NewRequest("GET", server.URL, nil)

	proxy.ReverseProxyHandler(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}
