#!/bin/sh

# script to execute `make $1` on all subfolder at the same time

stty -tostop

PIDS=""

if [ -f .env ]; then
    export $(cat .env | xargs)
fi

for folder in "${@:2}"; do
    make -C $folder $1 &
    PIDS="$! $PIDS"
done

killall() {
    kill $PIDS &> /dev/null
}

trap killall EXIT
trap killall SIGINT
trap killall SIGTERM

# -n to exit when any for the processes exits
wait -n $PIDS
