#!/bin/sh

make -C front dev &
FRONT_PID=$!
make -C game_server dev &
GAME_SERVER_PID=$!

function killall()
{
    kill $GAME_SERVER_PID
    kill $FRONT_PID
}

trap killall EXIT

wait $FRONT_PID
wait $GAME_SERVER_PID
