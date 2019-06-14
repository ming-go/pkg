package mhttp

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type server struct {
	chSd   chan struct{}
	chStop chan struct{}
	srv    *http.Server
}

type config struct {
	mux  *http.ServeMux
	addr string
	sync.Mutex
}

func (cfg *config) SetMux(mux *http.ServeMux) *config {
	cfg.Lock()
	cfg.mux = mux
	cfg.Unlock()
	return cfg
}

func (cfg *config) SetAddr(host string, port string) *config {
	cfg.Lock()
	cfg.addr = net.JoinHostPort(host, port)
	cfg.Unlock()
	return cfg
}

func NewConfig() *config {
	return &config{
		mux:  http.NewServeMux(),
		addr: ":8787",
	}
}

func NewDefaultServer(cfg *config) *server {
	return &server{
		chSd:   make(chan struct{}),
		chStop: make(chan struct{}),
		srv: &http.Server{
			Addr:         cfg.addr,
			Handler:      cfg.mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

func (srv *server) Run() error {
	go func() {
		select {
		case <-srv.chSd:
			if err := srv.srv.Shutdown(context.Background()); err != nil {
				log.Printf("HTTP server Shutdown: %v", err)
			}
		case <-srv.chStop:
			if err := srv.srv.Close(); err != nil {
				log.Printf("HTTP server Stop: %v", err)
			}
		}
	}()

	if err := srv.srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("HTTP server ListenAndServe: %v", err)
		return err
	}

	return nil
}

func (srv *server) Shutdown() {
	srv.chSd <- struct{}{}
}

func (srv *server) Stop() {
	srv.chStop <- struct{}{}
}
