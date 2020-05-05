package web

import (
	"flag"
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer() *Server {
	addr := *flag.String("address", ":1433", "address of server")
	flag.Parse()

	router := newRouter()

	return &Server{
		&http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func (ws *Server) Run() error {
	return ws.server.ListenAndServe()
}
