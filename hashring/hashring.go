/*
 * HashRing: a consistent hashing tracking structure
 */
package hashring

import (
	"crypto/sha1"
	"github.com/warlockcc/golibs/bisect"
	"strconv"
	"sync"
)

type Node interface {
	Address() string
}

// Key is the type of elements stored in the hashring
type Key [sha1.Size]byte

// a comparison function to be able to sort keys
func (k Key) Less(other bisect.Elem) bool {
	var o = other.(Key)
	for i := 0; i < len(k); i++ {
		if k[i] < o[i] {
			return true
		} else if o[i] < k[i] {
			return false
		}
	}
	return false // equal keys
}

type HashRing struct {
	num_replicas  uint
	sorted_keys   []bisect.Elem
	virtual_nodes map[Key]Node
	sync.Mutex
}

// a ring constructor
func New(replicas uint) *HashRing {
	return &HashRing{
		num_replicas:  replicas,
		sorted_keys:   make([]bisect.Elem, 0, 100),
		virtual_nodes: make(map[Key]Node),
		Mutex:         sync.Mutex{}}
}

// build a key from data
func (*HashRing) Hash(data []byte) (key Key) {
	var hasher = sha1.New()
	hasher.Write(data)
	hasher.Sum(key[:0])
	return
}

// add a suffix to the node address to simulate virtual nodes
func (r *HashRing) virtualNodeKey(n Node, i uint) Key {
	var replica = n.Address() + "_" + strconv.FormatUint(uint64(i), 10)
	return r.Hash([]byte(replica))
}

// Add nodes to the hashring
func (r *HashRing) AddNode(n Node) {
	for i := uint(0); i < r.num_replicas; i++ {
		var key = r.virtualNodeKey(n, i)
		if _, exists := r.virtual_nodes[key]; exists {
			panic("Hashring: two replicas hash equally")
		} else {
			r.Lock()
			r.sorted_keys = bisect.Insort(r.sorted_keys, key)
			r.virtual_nodes[key] = n
			r.Unlock()
		}
	}
}

// Remove node replicas from the ring
func (r *HashRing) RemoveNode(n Node) {
	for i := uint(0); i < r.num_replicas; i++ {
		var key = r.virtualNodeKey(n, i)
		if _, exists := r.virtual_nodes[key]; !exists {
			panic("Hashring: non existent node replica")
		} else {
			r.Lock()
			delete(r.virtual_nodes, key)
			r.sorted_keys = bisect.Remove(r.sorted_keys, key)
			r.Unlock()
		}
	}
}

// get the node where a key should be stored
// it should be stored in the rightmost node less than 'key'
func (r *HashRing) GetNode(key *Key) Node {
	r.Lock()
	defer r.Unlock()
	var replica_idx = bisect.Bisect(r.sorted_keys, key) - 1
	if replica_idx == -1 {
		replica_idx = len(r.sorted_keys) - 1
	}
	return r.virtual_nodes[r.sorted_keys[replica_idx].(Key)]
}

// get the 'num' nodes where a key can be stored replicated
// since the key is stored in the preceding virtual-node, replicas
// can be store in the n previous buckets in case the primary fails
func (r *HashRing) GetNodes(key *Key, num int) (list []Node) {
	list = make([]Node, 0, 3)
	r.Lock()
	var base_idx = bisect.Bisect(r.sorted_keys, key)
	for i := 1; i <= num; i++ {
		var idx = (base_idx - i + (num * len(r.sorted_keys))) % len(r.sorted_keys)
		list = append(list, r.virtual_nodes[r.sorted_keys[idx].(Key)])
	}
	r.Unlock()
	return
}

// get the right-open range [start, end) of keys a virtual-node can hold
func (r *HashRing) NodeKeys(n Node, replica uint) (start, end Key) {
	var key = r.virtualNodeKey(n, replica)
	r.Lock()
	var idx = bisect.Bisect(r.sorted_keys, key)
	end = r.sorted_keys[idx%len(r.sorted_keys)].(Key)
	if idx == 0 {
		start = r.sorted_keys[len(r.sorted_keys)-1].(Key)
	} else {
		start = r.sorted_keys[idx-1].(Key)
	}
	r.Unlock()
	return
}
