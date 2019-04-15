package backend

// Section naming convention: "<SECTION_NAME>"

// <GLOBAL>
const (
	// MPV (https://mpv.io) is a free, open source, and cross-platform media player
	MPV = "mpv"
	// VLC (https://www.videolan.org/vlc/index.html) is a free and open source
	// cross-platform multimedia player and framework that plays most multimedia
	// files as well as DVDs, Audio CDs, VCDs, and various streaming protocols.
	VLC = "vlc"
	OMX = "omx"
)

// <MPV>
const (
	mpvSocketName = "vscreen.sock"
)

// <VLC>
const (
	vlcPassword = "vscreen"
	vlcPort     = 8081
	vlcHost     = "localhost"
)

// <OMX>
const (
	omxBin = "omxplayer"
)

// <YTDL>
var (
	ytdlPath string
)
