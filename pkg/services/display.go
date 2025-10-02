package services

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type DisplayState string

const (
	DisplayStateOn      DisplayState = "On"
	DisplayStateOff     DisplayState = "Off"
	DisplayStateUnknown DisplayState = "Unknown"
)

type DisplayService struct {
	cmd        string
	getArgs    []string
	onArgs     []string
	offArgs    []string
	onMatch    *regexp.Regexp
	offMatch   *regexp.Regexp
	brightness uint8
}

type DisplayServiceOption func(*DisplayService)

func NewDisplayService(opts ...DisplayServiceOption) *DisplayService {
	svc := &DisplayService{
		cmd:        "vcgencmd",
		getArgs:    []string{"display_power"},
		onArgs:     []string{"display_power", "1"},
		offArgs:    []string{"display_power", "0"},
		onMatch:    regexp.MustCompile("display_power=1"),
		offMatch:   regexp.MustCompile("display_power=0"),
		brightness: 100,
	}

	for _, o := range opts {
		o(svc)
	}
	return svc
}

func WithDisplayCommand(cmd string) DisplayServiceOption {
	return func(s *DisplayService) {
		s.cmd = cmd
	}
}

func WithDisplayGetArgs(args ...string) DisplayServiceOption {
	return func(s *DisplayService) {
		s.getArgs = args
	}
}

func WithDisplayOnMatch(match *regexp.Regexp) DisplayServiceOption {
	return func(s *DisplayService) {
		s.onMatch = match
	}
}

func WithDisplayOffMatch(match *regexp.Regexp) DisplayServiceOption {
	return func(s *DisplayService) {
		s.offMatch = match
	}
}

func WithDisplayOnArgs(args ...string) DisplayServiceOption {
	return func(s *DisplayService) {
		s.onArgs = args
	}
}

func WithDisplayOffArgs(args ...string) DisplayServiceOption {
	return func(s *DisplayService) {
		s.offArgs = args
	}
}

func (r *DisplayService) GetState() (DisplayState, error) {
	cmd := exec.Command(r.cmd, r.getArgs...)
	var out strings.Builder
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return DisplayStateUnknown, err
	}
	if cmd.ProcessState.ExitCode() != 0 {
		return DisplayStateUnknown, fmt.Errorf("unexpected exit code %d", cmd.ProcessState.ExitCode())
	}
	s := out.String()
	if r.onMatch.MatchString(s) {
		return DisplayStateOn, nil
	} else if r.offMatch.MatchString(s) {
		return DisplayStateOff, nil
	}
	return DisplayStateUnknown, nil
}

func (r *DisplayService) SetState(state DisplayState) error {
	args := []string{}
	switch state {
	case DisplayStateOn:
		args = r.onArgs
	case DisplayStateOff:
		args = r.offArgs
	default:
		return fmt.Errorf("invalid state")
	}
	cmd := exec.Command(r.cmd, args...)
	if err := cmd.Run(); err != nil {
		return err
	}
	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf("unexpected exit code %d", cmd.ProcessState.ExitCode())
	}
	return nil
}

func (r *DisplayService) GetBrightness() uint8 {
	return r.brightness
}

func (r *DisplayService) SetBrightness(brightness uint8) {
	r.brightness = brightness
}
