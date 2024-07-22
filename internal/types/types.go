package types

import (
	"net/http"

	"github.com/hritesh04/thanos/internal/proxy"
)

type Server struct {
	Url   string
	Proxy proxy.IProxy
}

type IBalancer interface {
	Serve(w http.ResponseWriter, r *http.Request)
	AddServer(*Server)
}
