package mpv

type info struct{}

func (i *info) Title() string {
	return ""
}

func (i *info) Thumbnail() string {
	return ""
}

func (i *info) Volume() float64 {
	return 0.0
}

func (i *info) Position() float64 {
	return 0.0
}

func (i *info) State() string {
	return ""
}
