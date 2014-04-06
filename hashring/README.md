# hashring: a consistent hashing module for go

Hashring keeps a set of *nodes* which represent buckets which can hold keys. The
idea behind this structure is to keep rehashing (keys changing buckets) to a
minimum when adding or deleting nodes from the ring. The number of replicas lets
each node be split into several *virtual nodes* allowing the load to be spread
more uniformly and also to keep key handover fair on node arrival and departure.

## Example

		package main

		import (
			"fmt"
			"hashring"
		)

		type TestNode string

		func (n TestNode) Address() string {
			return string(n)
		}

		func main() {
			var ring = hashring.New(20)

			for i := 0; i < 50; i++ {
				var node = fmt.Sprintf("192.168.10.%v", 45+i)
				ring.AddNode(TestNode(node))
			}

			var data1 = []byte("Some data to store in some node of the ring")
			var data2 = []byte("Some more data that will probably be in other node")

			var data_key1 = ring.Hash(data1)
			var data_key2 = ring.Hash(data2)

			var storage_node_a = ring.GetNode(&data_key1)
			var storage_node_b = ring.GetNode(&data_key2)

			fmt.Printf("Data-1 should be stored at %v\n", storage_node_a.Address())
			fmt.Printf("Data-2 should be stored at %v\n", storage_node_b.Address())
		}


#### Sources
* http://michaelnielsen.org/blog/consistent-hashing/
* http://www.martinbroadhurst.com/Consistent-Hash-Ring.html
* http://www.tomkleinpeter.com/2008/03/17/programmers-toolbox-part-3-consistent-hashing/
* http://www.linuxjournal.com/article/6797?page=0,0
