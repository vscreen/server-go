package main

import (
	"log"

	"github.com/vscreen/server-go/server"
)

func main() {
	s, err := server.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(s.ListenAndServe(":8080"))
}
