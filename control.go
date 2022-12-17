package main

import (
	"log"
	"net"
	"net/http"
	"sync"
)

type ControlServer struct {
	SocketPath string
	Recorder   *Recorder
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func (c *ControlServer) ListenAndServe() error {
	listener, err := net.Listen("unix", c.SocketPath)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	var last []byte = nil
	var lock sync.Mutex

	mux.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		c.Recorder.Start()
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		buf := c.Recorder.Stop()
		lock.Lock()
		defer lock.Unlock()
		last = buf
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("/last", func(w http.ResponseWriter, r *http.Request) {
		lock.Lock()
		defer lock.Unlock()
		_, err := w.Write(last)
		if err != nil {
			log.Printf("http write error: %s", err)
		}
	})

	return http.Serve(listener, logRequest(mux))
}
