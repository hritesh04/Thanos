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
	var wg *sync.WaitGroup
	healthyServers := make(chan *types.Server, len(cfg.Servers))
	for _, backend := range cfg.Servers {
		wg.Add(1)
		go func(backend config.Server) {
			server := &types.Server{
				Url:            backend.Url,
				Proxy:          proxyFunc(backend.Url),
				HealthEndPoint: backend.HealthEndPoint,
			}
			if ok := leastConnection.CheckHostAlive(server.Url); ok {
				healthyServers <- server
			}
			wg.Done()
		}(backend)
	}
	wg.Wait()
	close(healthyServers)
	for server := range healthyServers {
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
	server.ActiveConnection++
	server.Proxy.ReverseProxyHandler(w, r)
	server.ActiveConnection--
}

func (lc *LeastConnection) Next() *types.Server {
	var leastConnectionIndex int
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	leastConnection := lc.servers[lc.i].ActiveConnection
	for i, server := range lc.servers {
		if server.ActiveConnection < leastConnection {
			leastConnectionIndex = i
		}
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
	lc.servers = append(lc.servers, proxyServer)
	lc.len++
}
