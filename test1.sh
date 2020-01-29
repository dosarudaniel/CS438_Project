#!/bin/bash


go build

xterm -title "Node 1" -hold -e "./CS438_Project -ipAddr=127.0.0.1:5000 -create -m=8"  &
xterm -title "Node 2" -hold -e "./CS438_Project -ipAddr=127.0.0.1:5001 -join -m=8 -existingNodeIp=127.0.0.1:5000 -existingNodeId=db"  &
xterm -title "Node 3" -hold -e "./CS438_Project -ipAddr=127.0.0.1:5002 -join -m=8 -existingNodeIp=127.0.0.1:5001 -existingNodeId=21"  &
