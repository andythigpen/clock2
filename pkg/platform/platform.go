package platform

const (
	ClockWidth   = 800
	Margin       = 20
	WindowWidth  = 1920
	WindowHeight = 480

	FPS = 30 // the Raspberry Pi supports at most 30FPS
)

const (
	PlatformPI = iota
	PlatformDesktop
)
