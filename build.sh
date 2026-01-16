#!/bin/sh

go build -o dist/paypal paypal/paypal.go
sudo setcap 'cap_net_raw,cap_net_admin=eip' dist/paypal
