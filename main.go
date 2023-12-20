package main

import (
	"glimmermloj/repository"
	"glimmermloj/router"
)

func main() {
	repository.Init()
	r := router.Init()
	certFile := "./brifin.top.crt"
	keyFile := "./brifin.top.key"
	err := r.RunTLS("0.0.0.0:40001", certFile, keyFile)
	if err != nil {
		return
	}
}
