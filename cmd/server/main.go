package main

import (
	"echo-test/app"
	"echo-test/config"
	"fmt"
	"log"
)

func main() {
	port := fmt.Sprintf(":%s", config.Port)

	e := app.EchoHandler()

	log.Fatal(e.Mux.Start(port))
}
