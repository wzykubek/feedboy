package main

import (
	"go.wzykubek.xyz/feedboy/pkg/server"
)

var directory = "./schemes"

func main() {
	srv := server.Server{
		Port:      "8080",
		SchemeDir: directory,
	}

	srv.LoadSchemes()

	srv.Start()
}
