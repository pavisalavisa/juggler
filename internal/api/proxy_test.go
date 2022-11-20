package api_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pavisalavisa/juggler/internal/api"
	"github.com/pavisalavisa/juggler/internal/proxy"
	"github.com/stretchr/testify/require"
)

const (
	validJsonReq     = `{"message":"This is a valid JSON"}`
	nonJsonReq       = `n^L<AD>in^M<95><DE>·<B6>e<A3>"T<D2>a<FA>b<BB>=<F0><E1>"<9B><B7>^<BA><DC>`
	mainService      = "api.main.com"
	secondaryService = "api.secondary.com"
)

func TestProxy_AnyMethod_ShouldCallBothServices(t *testing.T) {
	testCases := []struct {
		desc   string
		method string
		body   string
		target string
	}{
		{
			desc:   "Test proxy POST request should call both services",
			method: "POST",
			body:   validJsonReq,
			target: "/target",
		},
		{
			desc:   "Test proxy GET request should call both services",
			method: "GET",
			body:   validJsonReq,
			target: "/target",
		},
		{
			desc:   "Test proxy PUT request should call both services",
			method: "PUT",
			body:   validJsonReq,
			target: "/target",
		},
		{
			desc:   "Test proxy DELETE request should call both services",
			method: "DELETE",
			body:   validJsonReq,
			target: "/target",
		},
		{
			desc:   "Test proxy HEAD request should call both services",
			method: "HEAD",
			body:   validJsonReq,
			target: "/target",
		},
		{
			desc:   "Test proxy PATCH request should call both services",
			method: "PATCH",
			body:   validJsonReq,
			target: "/target",
		},
		{
			desc:   "Test proxy OPTIONS request should call both services",
			method: "OPTIONS",
			body:   validJsonReq,
			target: "/target",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// arrange
			r := httptest.NewRequest(tC.method, tC.target, strings.NewReader(tC.body))
			w := httptest.NewRecorder()
			c := &mockCaller{}
			p := fixtureProxy(c)

			// act
			p.ProxyHttpCall(w, r)

			// assert that both services were called

			calls := c.calls
			require.Len(t, calls, 2, "expected 2 calls to be performed")
			require.Equal(t, tC.target, c.calls[0].URL.Path)
			require.Equal(t, tC.target, c.calls[1].URL.Path)
		})
	}
}

func TestProxy_AnyRequestBody_ShouldProxyCalls(t *testing.T) {
	// arrange
	r := httptest.NewRequest("POST", "/any", strings.NewReader(nonJsonReq))
	w := httptest.NewRecorder()
	c := &mockCaller{}
	p := fixtureProxy(c)

	// act
	p.ProxyHttpCall(w, r)

	// assert that both services were called
	require.Len(t, c.calls, 2, "expected 2 calls to be performed")
	callBody, _ := c.calls[0].GetBody()
	bodyBytes, _ := io.ReadAll(callBody)

	require.Equal(t, []byte(nonJsonReq), bodyBytes, "first call should have unaltered body")

	callBody, _ = c.calls[1].GetBody()
	bodyBytes, _ = io.ReadAll(callBody)
	require.Equal(t, []byte(nonJsonReq), bodyBytes, "second call should have unaltered body")
}

func TestProxy_ShouldReturnMainServiceResponse(t *testing.T) {
	// arrange
	r := httptest.NewRequest("POST", "/any", strings.NewReader(nonJsonReq))
	w := httptest.NewRecorder()
	c := &mockCaller{}
	p := fixtureProxy(c)
	// act
	p.ProxyHttpCall(w, r)

	// assert
	mainSrvBody, _ := c.calls[0].GetBody()
	mainSrvRes, _ := io.ReadAll(mainSrvBody)

	proxyRes, err := io.ReadAll(w.Body)

	require.NoError(t, err, "reading proxy request should not return an error")

	require.Equal(t, proxyRes, mainSrvRes, "main service response should be returned by the proxy")
}

func TestProxy_ShouldReturnMainServiceHeaders(t *testing.T) {
	// arrange
	r := httptest.NewRequest("POST", "/any", strings.NewReader(nonJsonReq))
	w := httptest.NewRecorder()
	c := &mockCaller{}
	p := fixtureProxy(c)
	// act
	p.ProxyHttpCall(w, r)

	// assert
	mainSrvHeaders := c.calls[0].Header
	proxyHeaders := w.Result().Header

	require.Equal(t, mainSrvHeaders, proxyHeaders, "main service response headers should be returned by the proxy")
}

func TestProxy_MainServiceCallFail_ShouldReturnInternalServerError(t *testing.T) {
	// arrange
	r := httptest.NewRequest("POST", "/any", strings.NewReader(nonJsonReq))
	w := httptest.NewRecorder()
	c := &mockCaller{}
	p := fixtureProxy(c)

	c.onCall = func(r *http.Request) (<-chan *http.Response, <-chan error) {
		resCh, errCh := make(chan *http.Response), make(chan error)

		if r.URL.Host == mainService {
			errCh <- fmt.Errorf("Something went wrong calling the main service")

		}

		close(resCh)
		close(errCh)

		return resCh, errCh
	}

	// act
	p.ProxyHttpCall(w, r)

	// assert
	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode, "internal server error expected when calling the main service fails")
}

func TestProxy_SecondaryServiceFail_ShouldNotFail(t *testing.T) {

	// arrange
	r := httptest.NewRequest("POST", "/any", strings.NewReader(nonJsonReq))
	w := httptest.NewRecorder()
	c := &mockCaller{}
	p := fixtureProxy(c)

	c.onCall = func(r *http.Request) (<-chan *http.Response, <-chan error) {
		resCh, errCh := make(chan *http.Response), make(chan error)

		if r.URL.Host == secondaryService {
			errCh <- fmt.Errorf("Something went wrong calling the main service")

		}

		close(resCh)
		close(errCh)

		return resCh, errCh
	}

	// act
	p.ProxyHttpCall(w, r)

	// assert
	require.Equal(t, http.StatusOK, w.Result().StatusCode, "secondary service call failures should not be propagated")
}

func TestProxy_SecondaryServiceHanging_ShouldReturnMainResponseImmediatelly(t *testing.T) {

	// arrange
	r := httptest.NewRequest("POST", "/any", strings.NewReader(nonJsonReq))
	w := httptest.NewRecorder()
	c := &mockCaller{}
	p := fixtureProxy(c)

	c.onCall = func(r *http.Request) (<-chan *http.Response, <-chan error) {
		resCh, errCh := make(chan *http.Response), make(chan error)

		if r.URL.Host == secondaryService {
			errCh <- fmt.Errorf("Something went wrong calling the main service")

			time.Sleep(time.Second * 30)
		}

		close(resCh)
		close(errCh)

		return resCh, errCh
	}

	// act
	p.ProxyHttpCall(w, r)

	// assert
	require.Equal(t, http.StatusOK, w.Result().StatusCode, "main service response should be returned while secondary is hanging")
}

func fixtureProxy(caller proxy.Caller) *api.ProxyService {
	proxy := proxy.NewProxy(caller)
	return api.NewProxyService(proxy)
}

type mockCaller struct {
	proxy.NoOpCaller
	calls  []*http.Request
	onCall func(*http.Request) (<-chan *http.Response, <-chan error)
}

func (c *mockCaller) Call(req *http.Request) (<-chan *http.Response, <-chan error) {
	c.calls = append(c.calls, req)

	if c.onCall != nil {
		return c.onCall(req)
	}

	return c.NoOpCaller.Call(req)
}
