package http

import (
	"net/http"
	"time"
)

func IsHostAlive(url string) bool {
	tr := &http.Transport{
		IdleConnTimeout:       2 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	return resp.StatusCode == 200
}
