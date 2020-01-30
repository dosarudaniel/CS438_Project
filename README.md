# \[WIP\] Peerster v2.0

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
  ...
```

In order to give commands to your node, e.g., "look up this file" or
search "hello world.txt", you need to build the client. From the root
folder of the project run in terminal:
```
$ cd client 
$ go build
$ ./client --help
  ...
```

Thank you for taking time to check out our project!

daniel-florin.dosaru@epfl.ch
