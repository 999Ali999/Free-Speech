package main

import (
	"blogging/internal/initializers"
	"blogging/internal/server"
	"fmt"
)

func init() {
	initializers.LoadEnvs()
}

func main() {
	newServer := server.NewServer()

	err := newServer.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start newServer: %s", err))
	}
}
