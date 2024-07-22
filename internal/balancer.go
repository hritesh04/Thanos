package internal

import (
	"github.com/hritesh04/thanos/internal/algorithm"
	"github.com/hritesh04/thanos/internal/proxy"
	"github.com/hritesh04/thanos/internal/types"
	"github.com/hritesh04/thanos/pkg/config"
)

var balancer = map[string]func(config config.Config, proxyFunc proxy.ProxyFunc) types.IBalancer{
	"round-robin": algorithm.NewRoundRobin,
}

func NewLoadBalancer(cfg config.Config, proxyFunc proxy.ProxyFunc) types.IBalancer {
	return balancer[cfg.Algorithm](cfg, proxyFunc)
}
