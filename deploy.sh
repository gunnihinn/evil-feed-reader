#!/bin/bash

binary=evilfr
config=feeds.json

service=evil

# defines $user and $host
source secrets

rsync -avz "$binary" "$config" "$user@$host":
ssh "$user@$host" sudo systemctl restart "$service"
