package main

import (
	"resto-app-server/internal/repo"
	"resto-app-server/internal/router"
)

func main() {
	repository := repo.New()

	server := router.Init(repository)

	err := server.Run()
	if err != nil {
		panic(err)
	}
}
