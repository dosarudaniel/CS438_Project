#!/bin/bash


# Test the TransferFile function: send 2 chunks from sender to receiver "Hello, " and "World!"

go build

xterm -title "Node Sender 219 " -hold -e "./CS438_Project -name A -ipAddr=127.0.0.1:5000 -create -m 8 -v"  &
sleep 1
xterm -title "Node Receiver 33" -hold -e "./CS438_Project -name B -ipAddr=127.0.0.1:5001 -join -existingNodeIp 127.0.0.1:5000 -existingNodeId f98eeff24e2fced1a1336182a3e8775326262914cc4087066d9346431795ccdb -m 8 -v -role=2"  &
