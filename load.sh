#!/bin/bash
go build
export SERV=aws
./webBackend 2>>log.txt &
echo "Server is launched with PID:" 
echo $!
disown