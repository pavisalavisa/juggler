package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/hlog"
)

const MainServiceApi = "main.api.com"
const SecondaryServiceApi = "secondary.api.com"

type Proxy struct {
	caller Caller
}

func NewProxy(caller Caller) *Proxy {
	return &Proxy{caller: caller}
}

func (p Proxy) Orchestrate(req *http.Request) (*http.Response, error) {
	logger := hlog.FromRequest(req)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	subReqCtx, subCtxCancel := context.WithCancel(req.Context())

	mainReq := prepareRequest(body, req, subReqCtx)
	mainResCh, mainErrCh := p.caller.Call(mainReq)

	secondaryReq := prepareRequest(body, req, subReqCtx)
	secondaryResCh, secondaryErrCh := p.caller.Call(secondaryReq)

	for mainResCh != nil || mainErrCh != nil {
		select {
		case mainRes, ok := <-mainResCh:
			if !ok {
				logger.Debug().Msgf("Main response channel closed.")
				mainResCh = nil
				continue
			}

			logger.Debug().Msgf("Got response from the main service")
			go func() {
				defer subCtxCancel()
				compareResponses(mainRes, nil, secondaryResCh, secondaryErrCh)
			}()
			return mainRes, nil

		case mainErr, ok := <-mainErrCh:
			if !ok {
				logger.Debug().Msgf("Main response channel closed.")
				mainErrCh = nil
				continue
			}
			logger.Debug().Msgf("Got error from the main service")
			go func() {
				defer subCtxCancel()
				go compareResponses(nil, mainErr, secondaryResCh, secondaryErrCh)
			}()
			return nil, mainErr

		case <-req.Context().Done():
			logger.Warn().Msg("Orchestration halted because the orchestration is canceled")
			subCtxCancel()
			return nil, fmt.Errorf("Context was canceled")
		}
	}

	subCtxCancel()
	return nil, fmt.Errorf("Fatal failure, main caller didn't return any response.")
}

func prepareRequest(body []byte, original *http.Request, ctx context.Context) *http.Request {
	req := original.Clone(ctx)

	req.Host = MainServiceApi
	req.URL.Host = MainServiceApi
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	return req
}

func compareResponses(mainRes *http.Response, mainErr error, secondaryResCh <-chan *http.Response, secondaryErrCh <-chan error) {

}
