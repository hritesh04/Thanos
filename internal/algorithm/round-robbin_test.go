package algorithm

import (
	"testing"

	"github.com/hritesh04/thanos/internal/proxy"
	"github.com/hritesh04/thanos/pkg/config"
	"github.com/hritesh04/thanos/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func init() {
	logger.InitLogger()
}

func TestRoundRobbin(t *testing.T) {
	cfg := config.Config{
		Algorithm: "round-robin",
		Servers: []config.Server{
			{Url: "http://localhost:3001", HealthEndPoint: "/"},
			{Url: "http://localhost:3002", HealthEndPoint: "/"},
		},
	}
	proxy := func(url string) proxy.IProxy {
		return proxy.NewReverseProxy(url)
	}
	balancer := NewRoundRobin(cfg, proxy)

	assert.Nil(t, balancer)
}
