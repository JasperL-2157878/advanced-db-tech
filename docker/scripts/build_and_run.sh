#!/bin/sh

set -xe

cd source && go build -o ../build/server
cd .. && ./build/server