#!/bin/bash
go build
./webBackend 2>>log.txt &
echo "Server is launched with PID:" 
echo $!
disown