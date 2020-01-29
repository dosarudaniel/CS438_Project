#!/bin/bash


# Test the FindSuccessor latency
#
go build   # build the Peerster Node

## Create a N node Chord ring using m == 16
./CS438_Project -name A -ipAddr=127.0.0.1:5000 -create -m 16 -v &  # Genesis node

# Number of node in the Chord ring
N=4

for (( c = 1; c < $N; c++ ))
do
  e=""
  if [ $c -lt 10 ]; then
    e="0$c"
  else
    e=$c
  fi
  echo "Joining node 127.0.0.1:100$e"
  ./CS438_Project -name="B$e" -ipAddr="127.0.0.1:100$e" -join -existingNodeIp=127.0.0.1:5000 -existingNodeId="ccdb" -m 16 &> /dev/null &
done


## Wait for the finger tables to be computed
sleep 30

# Reset the output of the current test
echo "" > test4_out.txt

cd client
go build


for (( c = 0; c < 10; c++ ))
do
  sleep 0.2
  ./client -PeersterAddress=127.0.0.1:5000 -command=findSuccessor -ID="0000" &>> ../test4_out.txt
  ./client -PeersterAddress=127.0.0.1:5000 -command=findSuccessor -ID="4000" &>> ../test4_out.txt
  ./client -PeersterAddress=127.0.0.1:5000 -command=findSuccessor -ID="8000" &>> ../test4_out.txt
  ./client -PeersterAddress=127.0.0.1:5000 -command=findSuccessor -ID="c000" &>> ../test4_out.txt
done

cd ..


sleep 10

#
## Kill all CS438_Project processes
pkill -f CS438_Project
#sleep 1
#
## Compare with reference file
#diff test4_out.txt test4_ref.txt > test4_debug.txt
#diff_ret_code=$?
#
## Print result
#if [ $diff_ret_code == 0 ]; then
#  echo "TEST PASSED"
#  rm test4_out.txt  # Clean
#else
#  echo "TEST FAILED, see test4_*.txt files"
#fi