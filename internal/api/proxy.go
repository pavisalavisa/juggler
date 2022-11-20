package api

import (
	"net/http"

	"github.com/pavisalavisa/juggler/internal/proxy"
)

type ProxyService struct {
	orchestrator proxy.Proxy
}

func NewProxyService(orhcestrator *proxy.Proxy) *ProxyService {
	return &ProxyService{orchestrator: *orhcestrator}
}

func (p ProxyService) proxyHandler() http.Handler {
	return http.HandlerFunc(p.ProxyHttpCall)
}

func (p ProxyService) ProxyHttpCall(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented!")
}
