package player

type Property uint8

const (
	_ Property = iota
	PropPause
)

func mapPropName(property Property) string {
	var name string
	switch property {
	case PropPause:
		name = "pause"
	}

	return name
}

type Event struct {
	Name string
	Data interface{}
}

type EventHandler func(Event)
