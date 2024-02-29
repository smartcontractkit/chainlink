package mercury

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
)

type MercuryEndpointMock interface {
	URL() string
	Username() string
	Password() string
	CallCount() int
	RegisterHandler(http.HandlerFunc)
}

var _ MercuryEndpointMock = &SimulatedMercuryServer{}

var notFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
})

type SimulatedMercuryServer struct {
	server  *httptest.Server
	handler http.HandlerFunc

	callCounter atomic.Int32
}

func NewSimulatedMercuryServer() *SimulatedMercuryServer {
	srv := &SimulatedMercuryServer{
		handler:     notFoundHandler,
		callCounter: atomic.Int32{},
	}

	srv.server = httptest.NewUnstartedServer(srv)

	return srv
}

func (ms *SimulatedMercuryServer) URL() string {
	return ms.server.URL
}

func (ms *SimulatedMercuryServer) Username() string {
	return "username1"
}

func (ms *SimulatedMercuryServer) Password() string {
	return "password1"
}

func (ms *SimulatedMercuryServer) CallCount() int {
	return int(ms.callCounter.Load())
}

func (ms *SimulatedMercuryServer) RegisterHandler(h http.HandlerFunc) {
	ms.handler = h
}

func (ms *SimulatedMercuryServer) Start() {
	ms.server.Start()
}

func (ms *SimulatedMercuryServer) Stop() {
	ms.server.Close()
}

func (ms *SimulatedMercuryServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ms.callCounter.Add(1)
	ms.handler.ServeHTTP(w, r)
}
