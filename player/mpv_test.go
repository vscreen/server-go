package player

import "testing"

func TestMPV(t *testing.T) {
	mpv, err := MPVNew()
	if err != nil {
		t.Error(err)
	}
	defer mpv.Close()

	done := make(chan struct{})

	go mpv.Start()
	mpv.Add("https://www.youtube.com/watch?v=Dxf26RoPNTw")
	mpv.Stop()
	mpv.Add("https://www.youtube.com/watch?v=Kt4iMM2OigY")

	<-done
}
