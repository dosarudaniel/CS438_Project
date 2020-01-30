# Peerster v2.0: DHT based file sharing in Peerster

Peerster v2.0 is a peer-to-peer file sharing application based on
structured overlay network. It allows downloading files either by
knowing the file owner's ID and the filename or by a keywords-based
search (WIP).

For the implementation of the overlay network, we use
[Chord](https://pdos.csail.mit.edu/papers/chord:sigcomm01/chord_sigcomm.pdf).
Peerster v2.0 has now the following properties of chord:
 - Every peerster will store in its finger table at most `O(logN)` nodes, where N is the number
   of nodes in the network.
-  As well, every search for ID will take at most `O(logN)` "messages" (see `tests/test4.sh`)     

    
## Architecture    

In the picture below there are four Peersters. One of them is the genesis node, let's say A. Nodes B, C and D will join the ring network knowning one existing node. We can interact with a particular node using a client program or a web interface.
![Architecture](https://github.com/dosarudaniel/CS438_Project/blob/master/docs/Chord_ring_request_File.png) 

### Communication in Peerster
Since [Chord](https://pdos.csail.mit.edu/papers/chord:sigcomm01/chord_sigcomm.pdf) original algorithm uses Remote procedures calls, we decided to use [gRPC](https://grpc.io/), a high performance, general purpose Remote Procedure Call library. For Chord control messages this project uses Unary RPC, but for file download and upload we are using the Client/Server streaming RPC (see [gRPC concepts](https://grpc.io/docs/guides/concepts/) ). This library provides a simple interface and a fast development environment. Because the RPCs are using TCP, our Peersters will reliably transmit files between each other and it will also avoid congestions.

## Performance   
Our implementation is offering a `O(logN)` ID search time complexity in a Chord ring network, the graph below was obtained by runnning multiple (N) Peerster processes on the same localhost. Considering this, we can see a lower bound of just 0.5ms for the Find Successor query. Testing for various values of N, reveals the logarithmic curve below:
     
![Query time of findSuccessor RPC](https://github.com/dosarudaniel/CS438_Project/blob/dosarudaniel-improve-readme/docs/QueryTime_FindSuccessor.png)     

## Keyword-based search

We create a distributed hash table using our Chord overlay. Every node
will be responsible for keys, such that `hash(node.predecessor.IP) <
hash(key) <= hash(node.IP)`.

In order to support keyword-based search, we will be implementing the
algorithm sketched out in
[this paper](https://www.cs.utexas.edu/users/browne/CS395Tf2002/Papers/Keywordsearch.pdf).
We will build a (distributed) inverted index tree. In other words, we
map each keyword to a node in the DHT, which will store a list of
documents containing that keyword.

## How to build and run Peerster

In order to build the Peerster itself (as a node), you need to run `go
build` in terminal and run the program. Use `--help` to get information
what flags you need to provide to the Peerster:
```
$ go build
$ ./CS438_Project --help
  Usage of ./CS438_Project:
  -checkPredecessorInterval int
    	Number of seconds between two runs of CheckPredecessor Daemon (default 1)
  -create
    	Pass this flag to create a new Chord ring
  -existingNodeId string
    	The id to which this node should join
  -existingNodeIp string
    	ip:port for the existing Peerster in the Chord ring to join
  -fixFingerInterval int
    	Number of seconds between two runs of FixFingers Daemon (default 1)
  -ipAddr string
    	ip:port for the Peerster (default "127.0.0.1:5000")
  -join
    	Pass this flag to join to an existing Chord ring
  -m int
    	Number of bits in one node's id; max = 256, min = 4 (multiple of 4) (default 8)
  -name string
    	name of the Peerster
  -stabilizeInterval int
    	Number of seconds between two runs of Stabilize Daemon (default 1)
  -v	more verbosity of the program

```

In order to give commands to your node, e.g., "look up this file" or
search "hello world.txt", you need to build the client. From the root
folder of the project run in terminal:
```
$ cd client 
$ go build
$ ./client --help
  Usage of ./client:
  -ID string
    	Download: File owner's ID / FindSuccessor: ID for which the IP is requested 
  -PeersterAddress string
    	Peerster address to connect to
  -command string
    	Command to be sent to Peerster: download/upload/findSuccessor/search
  -file string
    	file name at owner
  -nameToStore string
    	Name used to store the downloaded file
  -query string
    	Search query (required for search command)
  -v	verbose mode
  -withDownload
    	Used with search command if you want to download one of the found results
```

Thank you for taking time to check out our project!

ulugbek.abdullaev@epfl.ch    
daniel-florin.dosaru@epfl.ch
