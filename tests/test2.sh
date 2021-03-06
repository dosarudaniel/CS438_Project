#!/bin/bash


# Test the FileSharing service and client interaction
cd ..

go build

# Removed current file1.txt for fair results
rm _download/file1.txt

xterm -title "Node Sender 219 " -hold -e "./CS438_Project -name A -ipAddr=127.0.0.1:5000 -create -m 8 -v"  &
sleep 1
cd client
go build
./client -PeersterAddress=127.0.0.1:5000 -command=download -file=file1.txt -nameToStore=file1.txt -ID=db -v
cd ..
echo "Diference between uploaded and downloaded:"
diff _upload/file1.txt _download/file1.txt

# Go back to tests
cd tests