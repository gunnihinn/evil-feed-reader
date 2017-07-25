#!/bin/bash

feeds=feeds.cfg
service=evil

# defines $user and $host
source remote/secrets

rsync -avz "$feeds" "$user@$host":
ssh "$user@$host" sudo service "$service" restart
