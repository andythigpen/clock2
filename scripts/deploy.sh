#!/bin/bash

set -e

/usr/local/go/bin/go build -tags drm,pi -o bin/clock
sudo /usr/bin/systemctl restart clock.service
