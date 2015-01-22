#!/bin/bash
go build
export SERV=aws
kill $(pgrep webBackend)
./webBackend 2>>log.txt &
echo "Server is launched with PID:" 
echo $!
disown