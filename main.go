package main

import (
	"log"

	"go.wzykubek.xyz/feedboy/pkg/parser"
)

func main() {
	scheme, err := parser.NewScheme("schemes/example.yml")
	if err != nil {
		panic(err)
	}

	parser := parser.Parser{Scheme: scheme}
	parser.Do()
	feed, err := parser.Feed.ToRss()
	if err != nil {
		log.Fatal(err)
	}

	println(feed)

}
