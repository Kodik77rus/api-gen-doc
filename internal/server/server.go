package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Kodik77rus/api-gen-doc/internal/config"
)

type Server struct {
	server *http.Server
}

func New(c *config.ServerConfig) (*Server, error) {
	server, err := configirateServer(c)
	if err != nil {
		return nil, err
	}

	return &Server{
		server: server,
	}, nil
}

func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func configirateServer(c *config.ServerConfig) (*http.Server, error) {
	rt, err := parceStingToInt(c.ReadTimeout)
	if err != nil {
		return nil, err
	}

	wrt, err := parceStingToInt(c.WriteTimeout)
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:         c.Port,
		ReadTimeout:  time.Duration(rt) * time.Second,
		WriteTimeout: time.Duration(wrt) * time.Second,
	}, nil
}

func parceStingToInt(str string) (int, error) {
	int, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return int, nil
}
