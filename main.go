package main

import (
	"echo-test/config"
	"fmt"
	"log"
)

func main() {
	port := fmt.Sprintf(":%s", config.Port)

	e := EchoHandler()

	log.Fatal(e.Start(port))
}
