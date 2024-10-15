package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsHostAlive(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)

	defer server.Close()

	assert.True(t, IsHostAlive(server.URL))
	assert.True(t, IsHostAlive("https://localhost:8080"))
}
