package player

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"

	"github.com/godbus/dbus"
	log "github.com/sirupsen/logrus"
)

func init() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	username := user.Username
	omxDBusAddr := fmt.Sprintf("/tmp/omxplayerdbus.%s", username)
	omxDBusPID := fmt.Sprintf("/tmp/omxplayerdbus.%s.pid", username)

	// initialize omx dbus daemon
	if err = exec.Command(omxBin).Run(); err != nil {
		log.Fatal(err)
	}

	dbusAddr, err := ioutil.ReadFile(omxDBusAddr)
	if err != nil {
		log.Fatal(err)
	}

	dbusPID, err := ioutil.ReadFile(omxDBusPID)
	if err != nil {
		log.Fatal(err)
	}

	os.Setenv("DBUS_SESSION_BUS_ADDRESS", string(dbusAddr))
	os.Setenv("DBUS_SESSION_BUS_PID", string(dbusPID))
}

type omxPlayer struct {
	conn *dbus.Conn
	obj  dbus.BusObject
}

func omxNew() (*omxPlayer, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	obj := conn.Object("org.mpris.MediaPlayer2.omxplayer", "/org/mpris/MediaPlayer2")

	return &omxPlayer{
		conn: conn,
		obj:  obj,
	}, nil
}

func (p *omxPlayer) play() error {
	return p.obj.Call("org.mpris.MediaPlayer2.Player.Play", 0).Err
}

func (p *omxPlayer) pause() error {
	return p.obj.Call("org.mpris.MediaPlayer2.Player.Pause", 0).Err
}

func (p *omxPlayer) add(url string) error {
	fmt.Println("adding", url)
	err := p.obj.Call("org.mpris.MediaPlayer2.Player.OpenUri", 0, url).Err
	if err != nil {
		fmt.Println("couldn't find omx session, starting one")
		return exec.Command(omxBin, "--no-osd", "--no-keys", url).Start()
	}
	return nil
}

func (p *omxPlayer) seek(position float64) error {
	dur, err := p.duration()
	if err != nil {
		return nil
	}

	microPos := int64(position * float64(dur))
	err = p.obj.Call("org.mpris.MediaPlayer2.Player.SetPosition", 0, dbus.ObjectPath("/not/used"), microPos).Err
	return nil
}

func (p *omxPlayer) close() {

}

func (p *omxPlayer) duration() (int64, error) {
	val, err := p.obj.GetProperty("org.mpris.MediaPlayer2.Player.Duration")
	if err != nil {
		return 0, err
	}

	dur := val.Value().(int64)
	return dur, nil
}
