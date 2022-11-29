package api

import (
	"io"
	"net/http"

	"github.com/pavisalavisa/juggler/internal/proxy"
	"github.com/rs/zerolog/hlog"
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
	logger := hlog.FromRequest(r)
	logger.Debug().Msg("Orchestrating request start")
	res, err := p.orchestrator.Orchestrate(r)

	ise := func(msg string, e error) {
		logger.Error().Err(e).Msg(msg)
		http.Error(w, "Unexpected issue proxying the request. Please try again later.", http.StatusInternalServerError)
	}

	if err != nil {
		ise("Something went wrong orchestrating the request", err)
		return
	}

	for k, vals := range res.Header {
		for _, v := range vals {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(res.StatusCode)

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		ise("Something went wrong reading response bytes", err)
		return
	}

	_, err = w.Write(resBytes)
	if err != nil {
		ise("Something went wrong writing response", err)
		return
	}
}
