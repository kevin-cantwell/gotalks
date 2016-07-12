#!/bin/bash

set -e

repo=$1
orighost=$2
port=$3

git clone https://$repo.git gotalk

cd gotalk

present -notes -http=0.0.0.0:$port -orighost=$orighost