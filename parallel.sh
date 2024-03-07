#!/bin/sh

# script to execute `make $1` on all subfolder at the same time

stty -tostop

PIDS=""

for folder in "${@:2}"; do
    make -C $folder $1 &
    PIDS="$! $PIDS"
done

killall() {
    kill $PIDS &> /dev/null
}

trap killall EXIT
trap killall SIGINT

# -n to exit when any for the processes exits
wait -n $PIDS
