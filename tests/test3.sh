#!/bin/bash

cd ..
# Test the FindSuccessor RPC

go build   # build the Peerster Node

./CS438_Project -name A -ipAddr=127.0.0.1:5000 -create -m 8 -v  &
./CS438_Project -name B -ipAddr=127.0.0.1:5001 -join -existingNodeIp 127.0.0.1:5000 -existingNodeId db -m 8 &> /dev/null &
./CS438_Project -name C -ipAddr=127.0.0.1:5003 -join -existingNodeIp 127.0.0.1:5000 -existingNodeId db -m 8 &> /dev/null &
./CS438_Project -name D -ipAddr=127.0.0.1:5004 -join -existingNodeIp 127.0.0.1:5000 -existingNodeId db -m 8 &> /dev/null &
./CS438_Project -name E -ipAddr=127.0.0.1:5005 -join -existingNodeIp 127.0.0.1:5000 -existingNodeId db -m 8 &> /dev/null &
./CS438_Project -name F -ipAddr=127.0.0.1:5006 -join -existingNodeIp 127.0.0.1:5000 -existingNodeId db -m 8 &> /dev/null &
./CS438_Project -name G -ipAddr=127.0.0.1:5007 -join -existingNodeIp 127.0.0.1:5000 -existingNodeId db -m 8 &> /dev/null &
./CS438_Project -name H -ipAddr=127.0.0.1:5008 -join -existingNodeIp 127.0.0.1:5000 -existingNodeId db -m 8 &> /dev/null &
./CS438_Project -name I -ipAddr=127.0.0.1:5009 -join -existingNodeIp 127.0.0.1:5000 -existingNodeId db -m 8 &> /dev/null &
./CS438_Project -name J -ipAddr=127.0.0.1:5010 -join -existingNodeIp 127.0.0.1:5000 -existingNodeId db -m 8 &> /dev/null &


# Wait for the finger tables to be computed
sleep 30

# Reset the output of the current test
echo "" > tests/test3_out.txt

cd client
go build
# Call findSuccessor for each node
# Request id for node B id = "21"
./client -PeersterAddress=127.0.0.1:5000 -command=findSuccessor -ID="21" > ../tests/test3_out.txt
sleep 0.5
# Request id for node C id = 9b
./client -PeersterAddress=127.0.0.1:5000 -command=findSuccessor -ID="9b" >> ../tests/test3_out.txt
sleep 0.5
# Request id for node D id = 83
./client -PeersterAddress=127.0.0.1:5000 -command=findSuccessor -ID="83" >> ../tests/test3_out.txt
sleep 0.5
# Request id for node E id = 5a
./client -PeersterAddress=127.0.0.1:5000 -command=findSuccessor -ID="5a" >> ../tests/test3_out.txt
sleep 0.5
# Request id for node F id = 68
./client -PeersterAddress=127.0.0.1:5000 -command=findSuccessor -ID="68" >> ../tests/test3_out.txt
sleep 0.5
# Request id for node G id = 0e
./client -PeersterAddress=127.0.0.1:5000 -command=findSuccessor -ID="0e" >> ../tests/test3_out.txt
sleep 0.5
# Request id for node H id = 70
./client -PeersterAddress=127.0.0.1:5004 -command=findSuccessor -ID="70" >> ../tests/test3_out.txt
sleep 0.5
# Request id for node I id = 03
./client -PeersterAddress=127.0.0.1:5003 -command=findSuccessor -ID="03" >> ../tests/test3_out.txt
sleep 0.5
# Request id for node J id = a8
./client -PeersterAddress=127.0.0.1:5001 -command=findSuccessor -ID="a8" >> ../tests/test3_out.txt
sleep 0.5
cd ..

# Kill all CS438_Project processes
pkill -f CS438_Project
sleep 1

# Compare with reference file
diff tests/test3_out.txt tests/test3_ref.txt > tests/test3_debug.txt
diff_ret_code=$?

# Print result
if [ $diff_ret_code == 0 ]; then
  echo "TEST PASSED"
  rm tests/test3_out.txt  # Clean
else
  echo "TEST FAILED, see tests/test3_*.txt files"
fi


# Go back to tests
cd tests