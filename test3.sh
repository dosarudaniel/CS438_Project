#!/bin/bash


# Test the FindSuccessor RPC

go build   # build the Peerster Node

xterm -title "Node A " -hold -e "./CS438_Project -name A -ipAddr=127.0.0.1:5000 -create -m 8 -v"  &
xterm -title "Node B " -hold -e "./CS438_Project -name B -ipAddr=127.0.0.1:5001 -join -existingNodeIp 127.0.0.1:5000 -existingNodeId f98eeff24e2fced1a1336182a3e8775326262914cc4087066d9346431795ccdb -m 8 -v -m 8 -v"  &
cd client
go build

sleep 10
# Request id for node B id =
./client -PeersterAddress=127.0.0.1:5000 -command=findSuccessor -ID="9724dfd73b0253a3b06ed53a5f9f1014997d5213d9cb1757363cbea588903321" -v
cd ..