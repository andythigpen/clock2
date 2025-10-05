//go:build !pi
// +build !pi

package screen

import "log/slog"

func ScreenOff() {
	slog.Warn("ScreenOff not supported on this platform")
}

func ScreenOn() {
	slog.Warn("ScreenOn not supported on this platform")
}

func IsScreenOn() (bool, error) {
	return true, nil
}
