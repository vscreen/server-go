package main

import "time"

func main() {
	p, _ := NewPlayer()
	time.Sleep(time.Second)
	p.Play()
	p.Pause()
	p.Stop()
	p.Add("POOP")
	p.Seek(30.0)

	p.Wait()
}
