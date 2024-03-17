package main

import (
	"log"
)

func main() {
	e := EchoHandler()

	log.Fatal(e.Start(":8080"))
}
