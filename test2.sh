#!/bin/bash


go build

xterm -title "Node 1 219 " -hold -e "./CS438_Project -name A -ipAddr=127.0.0.1:5000 -create -m 8 -v"  &
sleep 40
xterm -title "Node 2 33" -hold -e "./CS438_Project -name B -ipAddr=127.0.0.1:5001 -join -existingNodeIp 127.0.0.1:5000 -existingNodeId db -m 8 -v"  &
