package main

import (
	"os"
	"path/filepath"

	"go.wzykubek.xyz/feedboy/pkg/parser"
	"go.wzykubek.xyz/feedboy/pkg/server"
)

var directory = "./schemes"

func main() {
	srv := server.Server{
		Port:    "8080",
		Schemes: make(map[string]*parser.Scheme),
	}

	entries, err := os.ReadDir(directory)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		scheme, err := parser.NewScheme(filepath.Join(directory, entry.Name()))
		if err != nil {
			panic(err)
		}
		srv.Schemes[entry.Name()[:len(entry.Name())-4]] = scheme
	}

	srv.Start()
}
