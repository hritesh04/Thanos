package algorithm

import (
	"net/http"
	"sync"

	"github.com/hritesh04/thanos/internal/proxy"
	"github.com/hritesh04/thanos/internal/types"
	"github.com/hritesh04/thanos/pkg/config"
	"github.com/hritesh04/thanos/pkg/logger"
)

type RoundRobin struct {
	servers []*types.Server
	mutex   sync.Mutex
	// healthCheckInterval time.Duration
	len int
	i   int
}

func NewRoundRobin(cfg config.Config, proxyFunc proxy.ProxyFunc) types.IBalancer {
	roundRobin := &RoundRobin{}
	logger.Log.Info("Creating Round Robin Load Balancer")
	for _, backend := range cfg.Servers {
		server := &types.Server{
			Url:   backend,
			Proxy: proxyFunc(backend),
		}
		roundRobin.AddServer(server)
	}
	return roundRobin
}

func (rr *RoundRobin) Serve(w http.ResponseWriter, r *http.Request) {

	rr.Next().ReverseProxyHandler(w, r)
}

func (rr *RoundRobin) Next() proxy.IProxy {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()
	server := rr.servers[rr.i]
	rr.i = (rr.i + 1) % rr.len
	return server.Proxy
}

func (rr *RoundRobin) AddServer(proxyServer *types.Server) {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()
	rr.servers = append(rr.servers, proxyServer)
	rr.len++
	logger.Log.Info("Server Added", "Url", proxyServer.Url)
}
