package main

import (
	"github.com/physicist2018/url-shortener-go/internal/shortener"
)

func main() {
	if err := shortener.RunServer(); err != nil {
		panic(err)
	}
}
