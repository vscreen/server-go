package main

import (
	"fmt"
	"log"

	"github.com/vscreen/server-go/player"
)

func main() {
	p, err := player.New()
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	done := make(chan struct{})

	p.Handle(player.PropPause, func(e player.Event) {
		fmt.Println("handling", e.Name, e.Data)
	})

	go p.Start()
	p.Add("https://www.youtube.com/watch?v=b8eApABFWKE")
	p.Pause()

	<-done
}
