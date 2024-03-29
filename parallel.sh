#!/bin/sh

# script to execute `make $1` on all subfolder at the same time

stty -tostop

PIDS=""

make -C front $1 &
PIDS="$! $PIDS"
make -C creeps_server $1 &
PIDS="$! $PIDS"

killall() {
    kill $PIDS &> /dev/null
}

trap killall EXIT
trap killall SIGINT

# -n to exit when any for the processes exits
wait -n $PIDS
