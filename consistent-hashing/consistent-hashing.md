# Consistent Hashing in Go

This repository contains a Go implementation of **consistent hashing**, based on the concepts discussed in **Chapter 5** of the book *System Design Interview* by Alex Xu.

## Table of Contents

- [Introduction to Consistent Hashing](#introduction-to-consistent-hashing)
- [Problem Solved by Consistent Hashing](#problem-solved-by-consistent-hashing)
- [How Consistent Hashing Works](#how-consistent-hashing-works)
  - [Hash Ring](#hash-ring)
  - [Virtual Nodes](#virtual-nodes)
  - [Adding and Removing Nodes](#adding-and-removing-nodes)
- [Go Implementation](#go-implementation)
- [Example Usage](#example-usage)
- [References](#references)

## Introduction to Consistent Hashing
![consistent hasing](./'Screenshot 2024-09-27 at 1.14.02 PM.png')

Consistent hashing is a distributed hashing mechanism used in systems with dynamically changing nodes (servers). It ensures that the majority of keys are not reassigned when nodes join or leave the system, making it particularly useful in distributed systems like distributed caching, databases, and load balancing.

The main challenge addressed by consistent hashing is **minimizing disruption** when nodes are added or removed, while keeping key distribution relatively balanced across the nodes.

## Problem Solved by Consistent Hashing

In traditional hash-based systems, the removal or addition of a server can cause a significant number of keys to be reassigned, which can overload the system and lead to inefficiency. Specifically, if we use a modulo-based hashing scheme:

- **Hash(key) % N** where N is the number of nodes.
  
If a node is added or removed, most of the keys will need to be redistributed, which can cause instability in the system.

**Consistent Hashing** solves this problem by distributing both servers and keys on a circular hash space (hash ring). When a node is added or removed, only a small subset of keys are remapped, minimizing the impact on the system.

## How Consistent Hashing Works

### Hash Ring

In consistent hashing, both the servers (nodes) and keys are hashed into a circular hash space, often referred to as a "hash ring." The main property of this system is that the hash values are ordered in a circular fashion:

1. Servers are placed on the ring based on their hash value.
2. A key is hashed and mapped onto the ring.
3. The key is assigned to the first server node encountered in a clockwise direction.

### Virtual Nodes

To ensure more balanced key distribution, **virtual nodes** are used. A virtual node is a logical replica of a physical node placed at different points on the hash ring. These virtual nodes reduce the risk of uneven key distribution that might occur if servers are unevenly spaced on the hash ring.

For example, if a physical server `Server1` has 3 virtual nodes, it might be represented as `Server1#1`, `Server1#2`, `Server1#3`, each placed at a different position on the hash ring. This allows for better load balancing.

### Adding and Removing Nodes

- **Adding a Node**: When a new node is added, it is placed at several positions on the ring (due to virtual nodes). Keys that were originally mapped to the next node in a clockwise direction are reassigned to the new node.
  
- **Removing a Node**: When a node is removed, the keys that were mapped to it are reassigned to the next available node in the clockwise direction.

In both cases, the number of keys that need to be redistributed is proportional to the size of the node being added or removed, ensuring minimal disruption.

### Benefits

- **Minimal Key Movement**: Only a small portion of keys are reassigned when nodes are added or removed.
- **Scalability**: The system can easily scale by adding more nodes without drastically affecting the key distribution.
- **Load Balancing**: Virtual nodes ensure that load distribution remains even across the system.

## Go Implementation

The provided Go implementation includes:

- A **consistent hash ring** that distributes keys and nodes across a hash space.
- Support for **virtual nodes** to ensure even key distribution.
- Functions to **add or remove nodes** dynamically with minimal impact on the hash ring.
- Efficient lookup of nodes responsible for specific keys.

### Key Components

- **`hashKey(key string) uint32`**: Hashes a key using CRC32.
- **`AddNode(node string)`**: Adds a physical node and its virtual nodes to the hash ring.
- **`RemoveNode(node string)`**: Removes a node and its virtual nodes from the ring.
- **`GetNode(key string) string`**: Retrieves the node responsible for a given key by locating the closest node in the clockwise direction on the hash ring.
