package server

import (
	"fmt"
	"log"
	"net/http"

	"go.wzykubek.xyz/feedboy/pkg/parser"
)

type Server struct {
	schemes map[string]*parser.Scheme
	port    string
}

func NewServer(port string) *Server {
	return &Server{
		schemes: make(map[string]*parser.Scheme),
		port:    port,
	}
}

func (s *Server) AddScheme(schemeName string) error {
	scheme, err := parser.NewScheme("schemes/" + schemeName + ".yml") // TODO: Safe paths
	if err != nil {
		return fmt.Errorf("failed to load scheme: %w", err)
	}
	s.schemes[schemeName] = scheme

	return nil
}

func generateFeed(scheme *parser.Scheme) (string, error) {
	parser := parser.Parser{Scheme: scheme}
	parser.Do()
	feed, err := parser.Feed.ToRss()
	if err != nil {
		return "", fmt.Errorf("failed to generate feed: %w", err)
	}

	return feed, nil
}

func (s *Server) feedHandler(w http.ResponseWriter, r *http.Request) {
	schemeName := r.URL.Path[6:]

	if s.schemes[schemeName] == nil {
		http.NotFound(w, r)
		return
	}

	feed, err := generateFeed(s.schemes[schemeName])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, feed)
}

func (s *Server) Start() {
	http.HandleFunc("/feed/", s.feedHandler)
	log.Printf("Server starting on port %s", s.port)
	log.Printf("Configured feeds: %d", len(s.schemes))

	log.Fatal(http.ListenAndServe(":"+s.port, nil))
}
