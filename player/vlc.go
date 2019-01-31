package player

// Reference:
//  - https://wiki.videolan.org/VLC_HTTP_requests/
//  - https://github.com/videolan/vlc/blob/master/share/lua/http/requests/README.txt

import (
	"fmt"
	"net/http"
	netURL "net/url"
	"os"
	"os/exec"
)

type vlcPlayer struct {
}

func vlcNew() (*vlcPlayer, error) {
	cmd := exec.Command(
		vlcBin,
		"--fullscreen",
		"--no-loop",
		"--no-osd",
		"--play-and-stop",
		"--intf=macosx",
		"--extraintf=http",
		fmt.Sprintf("--http-password=%s", vlcPassword),
		fmt.Sprintf("--http-port=%d", vlcPort),
	)

	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &vlcPlayer{}, nil
}

func (p *vlcPlayer) close() {
}

func (p *vlcPlayer) get(args netURL.Values) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:%d/requests/status.xml?%s", vlcHost, vlcPort, args.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("", vlcPassword)

	return http.DefaultClient.Do(req)
}

func (p *vlcPlayer) play() error {
	args := netURL.Values{}
	args.Set("command", "pl_play")
	_, err := p.get(args)
	return err
}

func (p *vlcPlayer) pause() error {
	args := netURL.Values{}
	args.Set("command", "pl_pause")
	_, err := p.get(args)
	return err
}

func (p *vlcPlayer) add(url string) error {
	args := netURL.Values{}
	args.Set("command", "in_play")
	args.Set("input", url)
	_, err := p.get(args)
	return err
}

func (p *vlcPlayer) seek(position float64) error {
	args := netURL.Values{}
	args.Set("command", "seek")
	args.Set("val", fmt.Sprintf("%f%%", position*100))
	_, err := p.get(args)
	return err
}
