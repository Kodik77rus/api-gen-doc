package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/Kodik77rus/api-gen-doc/internal/config"
	handlers "github.com/Kodik77rus/api-gen-doc/internal/handelrs"
)

const apiPrefix = "/api"

type Server struct {
	server *http.Server
}

func New(c *config.ServerConfig) (*Server, error) {
	server, err := serverConfiguration(c)
	if err != nil {
		return nil, err
	}

	return &Server{
		server: server,
	}, nil
}

func (s *Server) Start(templateFolder string) error {
	router := httprouter.New()

	router.POST(apiPrefix+"/gendoc", handlers.GetGenDocHandler())
	router.POST(apiPrefix+"/find", handlers.FindDocs())

	router.ServeFiles(apiPrefix+"/download/*filepath", http.Dir(templateFolder))

	s.setRouter(router)

	if err := s.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) setRouter(router *httprouter.Router) {
	s.server.Handler = router
}

func serverConfiguration(c *config.ServerConfig) (*http.Server, error) {
	rt, err := parseStingToInt(c.ReadTimeout)
	if err != nil {
		return nil, err
	}

	wrt, err := parseStingToInt(c.WriteTimeout)
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:         c.Port,
		ReadTimeout:  time.Duration(rt) * time.Second,
		WriteTimeout: time.Duration(wrt) * time.Second,
	}, nil
}

func parseStingToInt(str string) (int, error) {
	int, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return int, nil
}
