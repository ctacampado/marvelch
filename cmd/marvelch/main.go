package main

import (
	"log"
	"marvel-chars/internal/marvelch"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	s := marvelch.MarvelCharSvc{}
	return s.Start("listening at ")
}
