package main

import (
	"flag"

	_ "github.com/physicist2018/url-shortener-go/internal/config"
	"github.com/physicist2018/url-shortener-go/internal/shortener"
)

func main() {
	flag.Parse()
	//fmt.Println(config.DefaultConfig)
	if err := shortener.RunServer(); err != nil {
		panic(err)
	}
}
