package core

import (
	"net/http"

	"github.com/hritesh04/thanos/internal/server"
	"github.com/hritesh04/thanos/pkg/config"
)

type LoadBalancer struct {
	servers      IServerPool
	algorithm    int
	healthchecks int
}

type IServerPool interface {
	GetRoundRobbin() *server.Server
	AddServer(string)
	// removeServer()
}

func (lb *LoadBalancer) handleRequest(w http.ResponseWriter, r *http.Request) {
	server := lb.getServer()
	server.Proxy.ServeHTTP(w, r)
}

func NewLoadBalancer(cfg config.Config) *LoadBalancer {
	return &LoadBalancer{
		servers:      new(server.ServerPool),
		algorithm:    0,
		healthchecks: 2,
	}
}

func (lb *LoadBalancer) AddServer(url string) {
	lb.servers.AddServer(url)
}

func (lb *LoadBalancer) getServer() *server.Server {
	if lb.algorithm == 0 {
		server := lb.servers.GetRoundRobbin()
		return server
	}
	return nil
}

func (lb *LoadBalancer) Start() {
	http.HandleFunc("/", lb.handleRequest)
	http.ListenAndServe(":3000", nil)
}
