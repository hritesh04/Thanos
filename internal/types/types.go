package types

import (
	"net/http"
	"sync"

	"github.com/hritesh04/thanos/internal/proxy"
)

type Server struct {
	Url              string
	ActiveConnection int
	Mutex            sync.Mutex
	Proxy            proxy.IProxy
}

type IBalancer interface {
	Serve(w http.ResponseWriter, r *http.Request)
	AddServer(*Server)
}
