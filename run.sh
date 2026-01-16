#!/bin/sh

./build.sh
sudo hciconfig hci0 down
dist/test
