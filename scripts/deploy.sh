#!/bin/bash

set -e

go build -tags drm,pi -o bin/clock
systemctl restart clock.service
