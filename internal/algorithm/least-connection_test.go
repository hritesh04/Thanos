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

func TestLeastConnection(t *testing.T) {
	cfg := config.Config{
		Algorithm: "least-connection",
		Servers: []config.Server{
			{
				Url: "http://locahost:3001", HealthEndPoint: "/",
			},
			{
				Url: "http://locahost:3001", HealthEndPoint: "/",
			},
		},
	}

	proxy := func(url string) proxy.IProxy {
		return proxy.NewReverseProxy(url)
	}

	balancer := NewLeastConnection(cfg, proxy)

	assert.Nil(t, balancer)
}
