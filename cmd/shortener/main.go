package main

import (
	"fmt"

	"github.com/physicist2018/url-shortener-go/internal/shortener"
)

func main() {
	var a map[string]string
	if _, ok := a["111"]; !ok {
		fmt.Println(ok)
	}
	if err := shortener.RunServer(); err != nil {
		panic(err)
	}
}
