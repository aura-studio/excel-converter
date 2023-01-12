package converter

type Mode int

var debugMode bool

const (
	ReleaseMode Mode = iota
	DebugMode
)

func SetMode(mode Mode) {
	if mode == DebugMode {
		debugMode = true
	}
}
