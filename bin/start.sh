#!/bin/bash

set -e

# Kill all child processes on exit
trap 'pkill -P $$' SIGINT SIGTERM EXIT

gotalk_host=$1
gotalk_port=$2

cd $GOPATH/src
present -notes -play=false -http=0.0.0.0:3999 -orighost=$gotalk_host >> /var/log/present.log 2>&1 &
PORT=$gotalk_port gotalk >> /var/log/gotalk.log 2>&1 &

tail -f /var/log/present.log /var/log/gotalk.log