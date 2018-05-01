#!/bin/bash

set -eu

binary=evil-feed-reader
man=evil-feed-reader.1
unit=evil-feed-reader.service
config=$HOME/.config/evil-feed-reader.yaml

service=evil-feed-reader

user=gmagnusson
host=droplet

rsync -az $binary $man $unit $config $user@$host:

ssh $user@$host sudo install $binary /usr/bin
ssh $user@$host sudo install $man /usr/share/man/man1
ssh $user@$host sudo install $unit /etc/systemd/system
ssh $user@$host sudo install evil-feed-reader.yaml /etc/xdg

ssh "$user@$host" sudo systemctl restart "$service"
