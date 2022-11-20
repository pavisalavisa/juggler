package proxy

import "net/http"

type Proxy struct {
	caller Caller
}

func NewProxy(caller Caller) *Proxy {
	return &Proxy{caller: caller}
}

func (p Proxy) Orchestrate(req *http.Request) (*http.Response, error) {
	panic("Not implemented!")
}
