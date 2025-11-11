package main

import "go.wzykubek.xyz/feedboy/pkg/server"

func main() {
	srv := server.NewServer("8080")
	srv.AddScheme("example")
	srv.Start()
}
