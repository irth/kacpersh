package main

import (
	"net"
	"net/http"
)

type ControlServer struct {
	SocketPath string
}

func (c *ControlServer) ListenAndServe() error {
	listener, err := net.Listen("unix", c.SocketPath)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	return http.Serve(listener, mux)
}
