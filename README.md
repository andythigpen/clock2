
Screen info:
```
Screen 0: minimum 320 x 200, current 1920 x 480, maximum 2048 x 2048
HDMI-1 connected primary 1920x480+0+0 left (normal left inverted right x axis y axis) 180mm x 270mm
   480x1920      59.42*+  53.09
```

Kernel cmdline options:
```
video=HDMI-A-1:480x1920@59.42e fbcon=rotate:3
```

Cross compiling (not working yet):
```
apt install gcc-aarch64-linux-gnu libgbm-dev
```

Building on a pi:
```
apt install libdrm-dev libegl1-mesa-dev libgles2-mesa-dev libgbm-dev
```

Splash screen on pi:
```
apt install rpd-plym-splash
```

In order to control the screen blanking on the pi, the `raylib-go` package was forked and two new functions were added:
```
extern int GetDrmConnectorPropertyValue(uint32_t property_id, uint64_t *property_value);
extern int SetDrmConnectorProperty(uint32_t property_id, uint64_t property_value);
```

To determine the property ids/values for DPMS, use
```
kmsprint -p
```

If for some reason they differ from the defaults `DPMS (2) = 0 (On) [On=0|Standby=1|Suspend=2|Off=3]`, the following env vars can be used to set new values:

```
DRM_DPMS_PROPERTY_ID
DRM_DPMS_PROPERTY_VALUE_ON
DRM_DPMS_PROPERTY_VALUE_OFF
```
