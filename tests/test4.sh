#!/bin/bash
# Test the FindSuccessor latency
#

usage(){
	echo "Usage: $0 N"  # N == number of nodes in the chord ring
	exit 1
}

# Safety checks
[[ $# -eq 0 ]] && usage

cd ..
go build   # build the Peerster Node

## Create a N node Chord ring using m == 16
./CS438_Project -name A -ipAddr=127.0.0.1:5000 -create -m 16 -v &  # Genesis node

# Number of node in the Chord ring
N=$1

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
echo "------   findSuccessor latency test (N = $1)   ------" > tests/test4_out.txt

cd client
go build


for (( c = 0; c < 10; c++ ))
do
  ./client -PeersterAddress=127.0.0.1:5000 -command=findSuccessor -ID="0000" -v >> ../tests/test4_out.txt
  ./client -PeersterAddress=127.0.0.1:5000 -command=findSuccessor -ID="4000" -v >> ../tests/test4_out.txt
  ./client -PeersterAddress=127.0.0.1:5000 -command=findSuccessor -ID="8000" -v >> ../tests/test4_out.txt
  ./client -PeersterAddress=127.0.0.1:5000 -command=findSuccessor -ID="c000" -v >> ../tests/test4_out.txt
done

cd ..


sleep 1
# Kill all CS438_Project processes
pkill -f CS438_Project

sleep 1

echo "------   findSuccessor latency test (N = $1)   ------" > ./tests/test4_stats_$1.txt
SUM1=0; for i in `cat ./tests/test4_out.txt | grep Time | grep 0000 | cut -d ":" -f2 | cut -d "," -f2`; do SUM1=$(($SUM1 + $i)); done; echo $((SUM1/10)) >> ./tests/test4_stats_$1.txt
SUM2=0; for i in `cat ./tests/test4_out.txt | grep Time | grep 4000 | cut -d ":" -f2 | cut -d "," -f2`; do SUM2=$(($SUM2 + $i)); done; echo $((SUM2/10)) >> ./tests/test4_stats_$1.txt
SUM3=0; for i in `cat ./tests/test4_out.txt | grep Time | grep 8000 | cut -d ":" -f2 | cut -d "," -f2`; do SUM3=$(($SUM3 + $i)); done; echo $((SUM3/10)) >> ./tests/test4_stats_$1.txt
SUM4=0; for i in `cat ./tests/test4_out.txt | grep Time | grep c000 | cut -d ":" -f2 | cut -d "," -f2`; do SUM4=$(($SUM4 + $i)); done; echo $((SUM4/10)) >> ./tests/test4_stats_$1.txt
echo "" >> ./tests/test4_stats_$1.txt
echo "Average: $(((SUM1+SUM2+SUM3+SUM4)/40))" >> ./tests/test4_stats_$1.txt

# Go back to tests
cd tests

