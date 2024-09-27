package main

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// Hash function using CRC32 algorithm
func hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

// ConsistentHash represents the consistent hashing ring
type ConsistentHash struct {
	hashCircle   map[uint32]string // Map from hash to node (server)
	sortedHashes []uint32          // Sorted list of hashes on the ring
	replicas     int               // Number of virtual nodes per physical node
	mutex        sync.RWMutex      // Read/Write lock for concurrent access
}

// NewConsistentHash creates a new ConsistentHash instance
func NewConsistentHash(replicas int) *ConsistentHash {
	return &ConsistentHash{
		hashCircle:   make(map[uint32]string),
		replicas:     replicas,
		sortedHashes: []uint32{},
	}
}

// AddNode adds a physical node (server) and its virtual nodes to the hash ring
func (ch *ConsistentHash) AddNode(node string) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	// Add virtual nodes for the given node
	for i := 0; i < ch.replicas; i++ {
		virtualNodeKey := node + strconv.Itoa(i)
		hash := hashKey(virtualNodeKey)
		ch.hashCircle[hash] = node
		ch.sortedHashes = append(ch.sortedHashes, hash)
	}

	// Sort hashes to maintain a sorted hash ring
	sort.Slice(ch.sortedHashes, func(i, j int) bool {
		return ch.sortedHashes[i] < ch.sortedHashes[j]
	})
}

// RemoveNode removes a physical node (server) and its virtual nodes from the hash ring
func (ch *ConsistentHash) RemoveNode(node string) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	// Remove virtual nodes for the given node
	for i := 0; i < ch.replicas; i++ {
		virtualNodeKey := node + strconv.Itoa(i)
		hash := hashKey(virtualNodeKey)
		delete(ch.hashCircle, hash)

		// Remove the hash from the sorted hash list
		index := sort.Search(len(ch.sortedHashes), func(i int) bool {
			return ch.sortedHashes[i] >= hash
		})

		if index < len(ch.sortedHashes) && ch.sortedHashes[index] == hash {
			ch.sortedHashes = append(ch.sortedHashes[:index], ch.sortedHashes[index+1:]...)
		}
	}
}

// GetNode finds the appropriate node for a given key
func (ch *ConsistentHash) GetNode(key string) string {
	ch.mutex.RLock()
	defer ch.mutex.RUnlock()

	if len(ch.sortedHashes) == 0 {
		return ""
	}

	// Hash the key
	hash := hashKey(key)

	// Binary search to find the first node hash that is >= hashed key
	idx := sort.Search(len(ch.sortedHashes), func(i int) bool {
		return ch.sortedHashes[i] >= hash
	})

	// If we reach the end of the slice, wrap around to the first node
	if idx == len(ch.sortedHashes) {
		idx = 0
	}

	// Return the node mapped to this hash
	return ch.hashCircle[ch.sortedHashes[idx]]
}

func main() {
	// Create a consistent hash ring with 3 virtual nodes per physical node
	hashRing := NewConsistentHash(3)

	// Add physical nodes (servers)
	hashRing.AddNode("Server1")
	hashRing.AddNode("Server2")
	hashRing.AddNode("Server3")

	// Get the responsible node for various keys
	keys := []string{"Key1", "Key2", "Key3", "Key4", "Key5"}

	for _, key := range keys {
		node := hashRing.GetNode(key)
		fmt.Printf("Key '%s' is handled by node '%s'\n", key, node)
	}

	// Remove a node and check redistribution of keys
	fmt.Println("\nRemoving Server2...")
	hashRing.RemoveNode("Server2")

	for _, key := range keys {
		node := hashRing.GetNode(key)
		fmt.Printf("Key '%s' is now handled by node '%s'\n", key, node)
	}
}

