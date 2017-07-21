#!/bin/bash

bin=evil-feed-reader
feeds=feeds.cfg
service=evil

# defines $user and $host
source ./secrets

rsync -avz "$bin" "$feeds" "$user@$host":
ssh "$user@$host" sudo service "$service" restart
