//go:build drm
// +build drm

package screen

import (
	"os"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const defaultPropertyId = 2

var (
	dpmsPropertyId       uint32 = 2
	dpmsPropertyValueOn  uint64 = 0
	dpmsPropertyValueOff uint64 = 3
)

func init() {
	e := os.Getenv("DRM_DPMS_PROPERTY_ID")
	if e != "" {
		id, err := strconv.ParseUint(e, 10, 32)
		if err == nil {
			dpmsPropertyId = uint32(id)
		}
	}
	e = os.Getenv("DRM_DPMS_PROPERTY_VALUE_ON")
	if e != "" {
		val, err := strconv.ParseUint(e, 10, 64)
		if err == nil {
			dpmsPropertyValueOn = val
		}
	}
	e = os.Getenv("DRM_DPMS_PROPERTY_VALUE_OFF")
	if e != "" {
		val, err := strconv.ParseUint(e, 10, 64)
		if err == nil {
			dpmsPropertyValueOn = val
		}
	}
}

func ScreenOff() {
	rl.SetDrmConnectorProperty(dpmsPropertyId, dpmsPropertyValueOff)
}

func ScreenOn() {
	rl.SetDrmConnectorProperty(dpmsPropertyId, dpmsPropertyValueOff)
}

func IsScreenOn() (bool, error) {
	val, err := rl.GetDrmConnectorPropertyValue(dpmsPropertyId)
	if err != nil {
		return false, err
	}
	return val == dpmsPropertyValueOn, nil
}
