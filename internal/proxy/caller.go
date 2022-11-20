package proxy

import "net/http"

type Caller interface {
	Call(req *http.Request) (<-chan *http.Response, <-chan error)
}

type NoOpCaller struct{}

func (c NoOpCaller) Call(req *http.Request) (<-chan *http.Response, <-chan error) {
	resCh, errCh := make(chan *http.Response), make(chan error)

	close(resCh)
	close(errCh)

	return resCh, errCh
}
