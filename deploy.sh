#!/bin/bash

binary=evilfr
config=evil.yaml

service=evil

# defines $user and $host
source secrets

rsync -avz "$binary" "$config" "$user@$host":
ssh "$user@$host" sudo service "$service" restart &
