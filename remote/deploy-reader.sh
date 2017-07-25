#!/bin/bash

bin=evil-feed-reader
service=evil

# defines $user and $host
source remote/secrets

rsync -avz "$bin" "$user@$host":
ssh "$user@$host" sudo service "$service" restart
