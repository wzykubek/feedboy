package server

import (
	"fmt"
	"log"
	"net/http"

	"go.wzykubek.xyz/feedboy/pkg/parser"
)

type Server struct {
	Port    string
	Schemes map[string]*parser.Scheme
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

	if s.Schemes[schemeName] == nil {
		http.NotFound(w, r)
		return
	}

	feed, err := generateFeed(s.Schemes[schemeName])
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
	log.Printf("Server starting on port %s", s.Port)
	log.Printf("Configured feeds: %d", len(s.Schemes))

	log.Fatal(http.ListenAndServe(":"+s.Port, nil))
}
