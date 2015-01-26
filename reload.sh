#!/bin/bash
go build
export SERV=aws
kill $(pgrep webBackend)
./webBackend 2>>err.txt 1>>out.txt &
echo "Server is launched with PID:" 
echo $!
disown