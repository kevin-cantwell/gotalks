#!/bin/bash

set -e

# Kill all child processes on exit
trap 'pkill -P $$' SIGINT SIGTERM EXIT

gotalks_host=$1
gotalks_port=$2

cd $GOPATH/src
present -notes -play=false -http=0.0.0.0:3999 -orighost=$gotalks_host >> /var/log/present.log 2>&1 &
PORT=$gotalks_port gotalks >> /var/log/gotalks.log 2>&1 &

tail -F /var/log/present.log /var/log/gotalks.log