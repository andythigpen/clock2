
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
