package algorithm

import (
	"net/http"
	"sync"

	"github.com/hritesh04/thanos/internal/proxy"
	"github.com/hritesh04/thanos/internal/types"
	"github.com/hritesh04/thanos/pkg/config"
)

type LeastConnection struct {
	servers []*types.Server
	mutex   sync.Mutex
	// healthCheckInterval time.Duration
	len int
	i   int
}

func NewLeastConnection(cfg config.Config, proxyFunc proxy.ProxyFunc) types.IBalancer {
	leastConnection := &LeastConnection{}
	for _, backend := range cfg.Servers {
		server := &types.Server{
			Url:   backend,
			Proxy: proxyFunc(backend),
		}
		leastConnection.AddServer(server)
	}
	return leastConnection
}

func (lc *LeastConnection) Serve(w http.ResponseWriter, r *http.Request) {
	server := lc.Next()
	server.Mutex.Lock()
	server.ActiveConnection++
	server.Mutex.Unlock()
	server.Proxy.ReverseProxyHandler(w, r)
	server.Mutex.Lock()
	server.ActiveConnection--
	server.Mutex.Unlock()
}

func (lc *LeastConnection) Next() *types.Server {
	var leastConnectionIndex int
	leastConnection := lc.servers[lc.i].ActiveConnection
	for i, server := range lc.servers {
		lc.mutex.Lock()
		if server.ActiveConnection < leastConnection {
			leastConnectionIndex = i
		}
		lc.mutex.Unlock()
	}
	lc.i = leastConnectionIndex
	return lc.servers[leastConnectionIndex]
}

func (lc *LeastConnection) AddServer(proxyServer *types.Server) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	lc.servers = append(lc.servers, proxyServer)
	lc.len++
}
