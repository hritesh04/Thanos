package proxy

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hritesh04/thanos/pkg/logger"
)

type IProxy interface {
	ReverseProxyHandler(http.ResponseWriter, *http.Request)
}

type ProxyClient struct {
	mutex *sync.Mutex
	proxy *httputil.ReverseProxy
}

type ProxyFunc func(string) IProxy

func NewReverseProxy(serverURL string) IProxy {
	// extract origin from serverURL
	origin, _ := url.Parse(serverURL)

	// create a director for httputil.ReverseProxy struct
	director := func(r *http.Request) {
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Header.Add("X-Origin-Host", origin.Host)
		r.URL.Scheme = "http"
		r.URL.Host = origin.Host
	}

	reverseProxy := &httputil.ReverseProxy{Director: director}

	return &ProxyClient{
		proxy: reverseProxy,
	}
}

func (p *ProxyClient) ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	// p.mutex.Lock()
	reqId := uuid.NewString()
	now := time.Now()
	r = r.WithContext(context.WithValue(r.Context(), "start", now))
	r = r.WithContext(context.WithValue(r.Context(), "reqId", reqId))
	// defer p.mutex.Unlock()
	p.proxy.ServeHTTP(w, r)
	startTime, ok := r.Context().Value("start").(time.Time)
	if !ok {
		log.Println("Error getting start time")
	}
	respId, ok := r.Context().Value("reqId").(string)
	if !ok {
		log.Println("Error getting request id")
	}
	diff := time.Since(startTime).String()
	logger.Log.Info("Completed Request "+respId, "response time", diff)
}
