package algorithm

import (
	"net/http"
	"sync"
	"time"

	"github.com/hritesh04/thanos/internal/proxy"
	"github.com/hritesh04/thanos/internal/types"
	"github.com/hritesh04/thanos/pkg/config"
	health "github.com/hritesh04/thanos/pkg/http"
	"github.com/hritesh04/thanos/pkg/logger"
)

type RoundRobin struct {
	servers             []*types.Server
	mutex               sync.Mutex
	healthCheckInterval time.Duration
	len                 int
	i                   int
}

func NewRoundRobin(cfg config.Config, proxyFunc proxy.ProxyFunc) types.IBalancer {
	roundRobin := &RoundRobin{}
	logger.Log.Info("Creating Round Robin Load Balancer")
	for _, backend := range cfg.Servers {
		server := &types.Server{
			Url:            backend.Url,
			Proxy:          proxyFunc(backend.Url),
			HealthEndPoint: backend.HealthEndPoint,
		}
		roundRobin.AddServer(server)
	}
	if len(roundRobin.servers) < 1 {
		logger.Log.Error("No healthy servers")
		return nil
	}
	return roundRobin
}

func (rr *RoundRobin) Serve(w http.ResponseWriter, r *http.Request) {
	rr.Next().ReverseProxyHandler(w, r)
}

func (rr *RoundRobin) CheckHostAlive(url string) bool {
	if alive := health.IsHostAlive(url); alive {
		logger.Log.Info("Server healthy", "url", url)
		return true
	}
	logger.Log.Error("Error checking server health", "url", url)
	return false
}

func (rr *RoundRobin) Next() proxy.IProxy {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()
	server := rr.servers[rr.i]
	rr.i = (rr.i + 1) % rr.len
	return server.Proxy
}

func (rr *RoundRobin) AddServer(proxyServer *types.Server) {
	alive := rr.CheckHostAlive(proxyServer.Url)
	if alive {
		rr.servers = append(rr.servers, proxyServer)
		rr.len++
		logger.Log.Info("Server Alive Added", "url", proxyServer.Url)
	} else {
		logger.Log.Info("Server Not Alive", "url", proxyServer.Url)
	}
}
