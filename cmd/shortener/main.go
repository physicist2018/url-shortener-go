package main

import (
	"github.com/physicist2018/url-shortener-go/internal/config"
	_ "github.com/physicist2018/url-shortener-go/internal/config"
	"github.com/physicist2018/url-shortener-go/internal/shortener"
)

func main() {
	_ = config.ConfigApp()
	//fmt.Println(config.DefaultConfig)
	if err := shortener.RunServer(); err != nil {
		panic(err)
	}
}
