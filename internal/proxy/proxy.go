package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type IProxy interface {
	ReverseProxyHandler(http.ResponseWriter, *http.Request)
}

type ProxyClient struct {
	proxy *httputil.ReverseProxy
}

type ProxyFunc func(string) IProxy

func NewReverseProxy(serverURL string) IProxy {
	//extract origin from serverURL
	origin, _ := url.Parse(serverURL)

	//create a director for httputil.ReverseProxy struct
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
	p.proxy.ServeHTTP(w, r)
}
