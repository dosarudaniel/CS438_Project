#!/bin/bash


# Test the FileSharing service and Client interaction

go build

rm _Download/file1.txt

xterm -title "Node Sender 219 " -hold -e "./CS438_Project -name A -ipAddr=127.0.0.1:5000 -create -m 8 -v"  &
sleep 1
cd Client
go build
./Client -PeersterAddress=127.0.0.1:5000 -file=file1.txt -ownersID=f98eeff24e2fced1a1336182a3e8775326262914cc4087066d9346431795ccdb -v
cd ..
echo "Diference between uploaded and downloaded:"
diff _Upload/file1.txt _Download/file1.txt