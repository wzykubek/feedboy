package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"go.wzykubek.xyz/feedboy/pkg/parser"
)

type Server struct {
	Port      string
	SchemeDir string
	schemes   map[string]*parser.Scheme
}

func (s *Server) LoadSchemes() error {
	s.schemes = make(map[string]*parser.Scheme)

	entries, err := os.ReadDir(s.SchemeDir)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		scheme, err := parser.NewScheme(filepath.Join(s.SchemeDir, entry.Name()))
		if err != nil {
			panic(err)
		}
		s.schemes[entry.Name()[:len(entry.Name())-4]] = scheme
	}

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

	err := s.LoadSchemes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	log.Printf("Server starting on port %s", s.Port)
	log.Printf("Configured feeds: %d", len(s.schemes))

	log.Fatal(http.ListenAndServe(":"+s.Port, nil))
}
