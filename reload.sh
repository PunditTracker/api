#!/bin/bash
go build
export SERV=local
export AWS_ACCESS_KEY_ID="AKIAJUAA7V5XHKH2A75Q"
export AWS_SECRET_ACCESS_KEY="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"
kill $(pgrep webBackend)
./webBackend 2>>err.txt 1>>out.txt &
echo "Server is launched with PID:" 
echo $!
disown