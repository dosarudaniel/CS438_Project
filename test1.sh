#!/bin/bash


go build

xterm -title "Node 1" -hold -e "./CS438_Project -ipAddr=127.0.0.1:5000 -create -m=8"  &
xterm -title "Node 2" -hold -e "./CS438_Project -ipAddr=127.0.0.1:5001 -join -m=8 -existingNodeIp=127.0.0.1:5000 -existingNodeId=f98eeff24e2fced1a1336182a3e8775326262914cc4087066d9346431795ccdb"  &
xterm -title "Node 3" -hold -e "./CS438_Project -ipAddr=127.0.0.1:5002 -join -m=8 -existingNodeIp=127.0.0.1:5001 -existingNodeId=9724dfd73b0253a3b06ed53a5f9f1014997d5213d9cb1757363cbea588903321"  &
