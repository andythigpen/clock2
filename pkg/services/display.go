package services

type DisplayState string

const (
	DisplayStateOn  DisplayState = "On"
	DisplayStateOff DisplayState = "Off"
)

type DisplayService struct {
	brightness uint8        // the desired brightness
	state      DisplayState // the desired state (may not be the actual state)
}

func NewDisplayService() *DisplayService {
	svc := &DisplayService{
		brightness: 100,
		state:      DisplayStateOn,
	}
	return svc
}

func (r *DisplayService) GetState() DisplayState {
	return r.state
}

func (r *DisplayService) SetState(state DisplayState) {
	r.state = state
}

func (r *DisplayService) GetBrightness() uint8 {
	return r.brightness
}

func (r *DisplayService) SetBrightness(brightness uint8) {
	r.brightness = brightness
}
