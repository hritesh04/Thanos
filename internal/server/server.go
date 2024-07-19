package server

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	url    string
	Proxy  *httputil.ReverseProxy
	Mutex  sync.Mutex
	weight float32
}

type ServerPool struct {
	servers []*Server
	mutex   sync.Mutex
}

func (s *ServerPool) GetRoundRobbin() *Server {
	return s.servers[0]
}

func (s *ServerPool) AddServer(url string) {
	server := new(Server)
	server.url = url
	server.Proxy = s.createReverseProxy(url)
	s.mutex.Lock()
	s.servers = append(s.servers, server)
	s.mutex.Unlock()
	log.Printf("Added new server to the server pool %s", url)
}

func (s *ServerPool) createReverseProxy(serverURL string) *httputil.ReverseProxy {
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

	return reverseProxy
}
