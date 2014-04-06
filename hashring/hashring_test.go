package hashring

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type TestNode string

func (n *TestNode) Address() string {
	return string(*n)
}

// test the hashring's hash function
func TestHash(t *testing.T) {
	var r = New(1)
	var sha1strings = map[string]string{
		"just a test string":         "3f0cf2e3d9e5903e839417dfc47fed6bfa6457f6",
		"another hashed string":      "740e0848eb95780f58b1fc973966b942bf245a14",
		"this last one and thats it": "9225afc5ba1cd274927cfc367d0c2b188f2f13a3"}

	for phrase, hash := range sha1strings {
		var hr_hash = r.Hash([]byte(phrase))
		if hex.EncodeToString(hr_hash[:]) != hash {
			t.FailNow()
		}
	}
}

// test our comparison function behaves correctly
func TestKeyCmp(t *testing.T) {
	var keys [10]Key
	copy(keys[0][:], bytes.NewBufferString("00000000000000000000").Bytes())
	copy(keys[1][:], bytes.NewBufferString("000000000aaaaaaa2222").Bytes())
	copy(keys[2][:], bytes.NewBufferString("0123456789abcdef0123").Bytes())
	copy(keys[3][:], bytes.NewBufferString("09854329584sbfecd234").Bytes())
	copy(keys[4][:], bytes.NewBufferString("09cebda4321cdef01234").Bytes())
	copy(keys[5][:], bytes.NewBufferString("123446789abcdef01234").Bytes())
	copy(keys[6][:], bytes.NewBufferString("123456789abcdef01234").Bytes())
	copy(keys[7][:], bytes.NewBufferString("bbbccaefd345ba984214").Bytes())
	copy(keys[8][:], bytes.NewBufferString("f0f0f0f0f0f0f0f0f0f0").Bytes())
	copy(keys[9][:], bytes.NewBufferString("ffffffffffffffffffff").Bytes())

	for i := 1; i < len(keys); i++ {
		if !keys[i-1].Less(keys[i]) || keys[i].Less(keys[i-1]) || keys[i].Less(keys[i]) {
			t.FailNow()
		}
	}
}

func nodeFromTemplate(nodenum int) *TestNode {
	var n = TestNode(fmt.Sprintf("testnode-%d", nodenum))
	return &n
}

// concurrent test add/remove nodes
func TestAddRemove(t *testing.T) {
	var wg sync.WaitGroup
	var num_nodes = 500
	var r = New(10)
	// test adding some nodes
	for i := 0; i < num_nodes; i++ {
		wg.Add(1)
		go func(j int) {
			r.AddNode(nodeFromTemplate(j))
			wg.Done()
		}(i)
	}
	wg.Wait()
	if len(r.virtual_nodes) != num_nodes*int(r.num_replicas) {
		t.FailNow()
	}
	// test deleting some nodes
	var removed = 0
	for i := rand.Intn(5); i < num_nodes; i += 1 + rand.Intn(7) {
		wg.Add(1)
		removed++
		go func(j int) {
			r.RemoveNode(nodeFromTemplate(j))
			wg.Done()
		}(i)
	}
	wg.Wait()
	if len(r.virtual_nodes) != (num_nodes-removed)*int(r.num_replicas) {
		t.FailNow()
	}
}

// a random key generator
func randomKey() (key Key) {
	for i := 0; i < len(key); i++ {
		key[i] = byte(rand.Intn(256))
	}
	return
}

// test retriving nodes
func TestGetNode(t *testing.T) {
	var wg sync.WaitGroup
	var num_nodes = 500
	var r = New(20)
	for i := 0; i < num_nodes; i++ {
		wg.Add(1)
		go func(j int) {
			r.AddNode(nodeFromTemplate(j))
			wg.Done()
		}(i)
	}
	wg.Wait()

	for n := 0; n < 5000; n++ {
		var key = randomKey()
		var node = r.GetNode(&key)
		var found_replicas = 0
		for i := uint(0); i < r.num_replicas; i++ {
			var start, end = r.NodeKeys(node, i)
			if !key.Less(start) && key.Less(end) ||
				 // consider the case when the key fits between the last and the first
				 start == r.sorted_keys[len(r.sorted_keys)-1] && end == r.sorted_keys[0] &&
				 (!key.Less(start) != key.Less(end)) {
				found_replicas++
			}
		}
		if found_replicas != 1 {
			t.Fatalf("Found %v replicas\n", found_replicas)
		}
	}
}


func ExampleHashring() {
	var ring = New(20)

	for i := 0; i < 50; i++ {
		var node = TestNode(fmt.Sprintf("192.168.10.%v", 45+i))
		ring.AddNode(&node)
	}

	var data1 = []byte("Some data to store in some node of the ring")
	var data2 = []byte("Some more data that will probably be in other node")

	var data_key1 = ring.Hash(data1)
	var data_key2 = ring.Hash(data2)

	var storage_node_a = ring.GetNode(&data_key1)
	var storage_node_b = ring.GetNode(&data_key2)

	fmt.Printf("Data-1 should be stored at %v\n", storage_node_a.Address())
	fmt.Printf("Data-2 should be stored at %v\n", storage_node_b.Address())

	// Output:
	// Data-1 should be stored at 192.168.10.89
	// Data-2 should be stored at 192.168.10.61
}
