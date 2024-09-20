package algorithm

import (
	"net/http"
	"sync"

	"github.com/hritesh04/thanos/internal/proxy"
	"github.com/hritesh04/thanos/internal/types"
	"github.com/hritesh04/thanos/pkg/config"
	health "github.com/hritesh04/thanos/pkg/http"
	"github.com/hritesh04/thanos/pkg/logger"
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
			Url:            backend.Url,
			Proxy:          proxyFunc(backend.Url),
			HealthEndPoint: backend.HealthEndPoint,
		}
		leastConnection.AddServer(server)
	}
	if len(leastConnection.servers) < 1 {
		logger.Log.Error("No healthy servers")
		return nil
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

func (lc *LeastConnection) CheckHostAlive(url string) bool {
	if alive := health.IsHostAlive(url); alive {
		logger.Log.Info("Server healthy", "url", url)
		return true
	}
	logger.Log.Error("Error checking server health", "url", url)
	return false
}

func (lc *LeastConnection) AddServer(proxyServer *types.Server) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	lc.servers = append(lc.servers, proxyServer)
	lc.len++
}
