#!/bin/bash


# Test the FileSharing service and client interaction
cd ..

go build

# Removed current file1.txt for fair results
rm _download/file1.txt

xterm -title "Node A " -hold -e "./CS438_Project -name A -ipAddr=127.0.0.1:5000 -create -m 8" &
xterm -title "Node B " -hold -e "./CS438_Project -name B -ipAddr=127.0.0.1:10001 -join -existingNodeIp 127.0.0.1:5000 -existingNodeId db -m 8" &
xterm -title "Node C " -hold -e "./CS438_Project -name C -ipAddr=127.0.0.1:10002 -join -existingNodeIp 127.0.0.1:5000 -existingNodeId db -m 8 -v" &
xterm -title "Node D " -hold -e "./CS438_Project -name D -ipAddr=127.0.0.1:10003 -join -existingNodeIp 127.0.0.1:5000 -existingNodeId db -m 8" &
xterm -title "Node E " -hold -e "./CS438_Project -name E -ipAddr=127.0.0.1:10004 -join -existingNodeIp 127.0.0.1:5000 -existingNodeId db -m 8" &

sleep 15
cd client
go build
./client -PeersterAddress=127.0.0.1:10004 -command=download -file=file1.txt -nameToStore=file1.txt -ID="db"
cd ..
echo "Diference between uploaded and downloaded:"
diff _upload/file1.txt _download/file1.txt

# Go back to tests
cd tests