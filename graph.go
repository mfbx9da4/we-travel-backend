package main

import (
	"container/heap"
	"fmt"
	"math"

	//"math"
	"sync"
	//"github.com/cheekybits/genny/generic"
)

var exists = struct{}{}

// Node a single node that composes the tree
type Node struct {
	Value Coordinate
	Hash  string
}

// GetOrCreateNode for create node from coords
func GetOrCreateNode(coords Coordinate) Node {
	hash := HashCoordinate(coords)
	if graph.nodes[hash] == nil {
		return Node{coords, hash}
	}
	return *graph.nodes[hash]
}

func (n *Node) String() string {
	return fmt.Sprintf("%v", n.Value)
}

// SetOfNodes Set of nodes
type SetOfNodes = map[*Node]struct{}

// Graph Basic graph complete with concurrency safe lock
type Graph struct {
	nodes map[string]*Node
	edges map[Node]SetOfNodes
	lock  sync.RWMutex
}

// AddNode adds a node to the graph
func (graph *Graph) AddNode(n *Node) {
	graph.lock.Lock()
	if graph.nodes == nil {
		graph.nodes = make(map[string]*Node)
	}
	graph.nodes[n.Hash] = n
	graph.lock.Unlock()
}

// AddEdge adds an edge to the graph
func (graph *Graph) AddEdge(n1, n2 *Node) {
	graph.lock.Lock()
	if graph.edges == nil {
		graph.edges = make(map[Node]SetOfNodes)
	}
	if graph.edges[*n1] == nil {
		graph.edges[*n1] = make(SetOfNodes)
	}
	if graph.edges[*n2] == nil {
		graph.edges[*n2] = make(SetOfNodes)
	}
	graph.edges[*n1][n2] = exists
	graph.edges[*n2][n1] = exists
	graph.lock.Unlock()
}

// Print graph
func (graph *Graph) String() {
	graph.lock.RLock()
	s := ""
	// for i := 0; i < len(graph.nodes); i++ {
	for _, node := range graph.nodes {
		s += node.String() + " -> "
		near := graph.edges[*node]
		for j := range near {
			s += j.String() + " "
		}
		s += "\n"
	}
	fmt.Println(s)
	graph.lock.RUnlock()
}

type QueueItemValue struct {
	Node     Node
	Path     []Node
	Distance float64
}

type NodeQueue struct {
	items []QueueItemValue
	lock  sync.RWMutex
}

// New creates a new NodeQueue
func (s *NodeQueue) New() *NodeQueue {
	s.lock.Lock()
	s.items = []QueueItemValue{}
	s.lock.Unlock()
	return s
}

// Enqueue adds an Node to the end of the queue
func (s *NodeQueue) Enqueue(t QueueItemValue) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

// Dequeue removes an Node from the start of the queue
func (s *NodeQueue) Dequeue() *QueueItemValue {
	s.lock.Lock()
	item := s.items[0]
	s.items = s.items[1:len(s.items)]
	s.lock.Unlock()
	return &item
}

// Front returns the item next in the queue, without removing it
func (s *NodeQueue) Front() *QueueItemValue {
	s.lock.RLock()
	item := s.items[0]
	s.lock.RUnlock()
	return &item
}

// IsEmpty returns true if the queue is empty
func (s *NodeQueue) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items) == 0
}

// Size returns the number of Nodes in the queue
func (s *NodeQueue) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}

// FindNode find node from graph
func (graph *Graph) FindNode(coords Coordinate) *Node {
	hash := HashCoordinate(coords)
	var node = graph.nodes[hash]
	if node == nil {
		var minDistance float64 = math.MaxFloat64
		// Probably there is a better algo for this, just doing the brute force sorry :(
		for _, cur := range graph.nodes {
			value := cur.Value
			dx := coords[0] - value[0]
			dy := coords[1] - value[1]
			distance := math.Sqrt(dx*dx + dy*dy)
			if distance < minDistance {
				node = cur
				minDistance = distance
			}
		}
	}
	return node
}

// Route Best route and distance of route
type Route struct {
	Path     []Node
	Distance float64
}

/*
FindPath Uses A* routing to find shortest path
(Keeps a min heap sorted by elapsed + remaing distance).
*/
func (graph *Graph) FindPath(src, dest *Node) Route {
	graph.lock.RLock()

	// Init priority queue
	var pqueue = make(PriorityQueue, 1)
	var rootPath = []Node{*src}
	var rootValue = QueueItemValue{*src, rootPath, 0}
	pqueue[0] = &QueueItem{
		Value:    &rootValue,
		Priority: 0,
		Index:    0,
	}
	heap.Init(&pqueue)

	// Keep track of visited
	visited := make(map[*Node]bool)
	for {
		if pqueue.Len() == 0 {
			break
		}
		pqitem := pqueue.Pop().(*QueueItem)
		cur := pqitem.Value
		node := cur.Node
		visited[&node] = true
		children := graph.edges[node]

		for child := range children {
			dx := (node.Value[0] - child.Value[0])
			dy := (node.Value[1] - child.Value[1])
			remaingDx := (dest.Value[0] - child.Value[0])
			remainingDy := (dest.Value[1] - child.Value[1])
			elapsed := math.Sqrt(dx*dx+dy*dy) + cur.Distance
			remaining := math.Sqrt(remaingDx*remaingDx + remainingDy*remainingDy)

			if *child == *dest {
				fmt.Println(*child, *dest)
				path := append(cur.Path, *child)
				return Route{path, elapsed}
			}

			if !visited[child] {
				// TODO: Only add to path if different gradient
				path := append(cur.Path, *child)
				queueItem := QueueItemValue{*child, path, elapsed}
				newItem := QueueItem{
					Value:    &queueItem,
					Priority: elapsed + remaining,
				}
				heap.Push(&pqueue, &newItem)
				visited[child] = true
			}
		}
	}
	graph.lock.RUnlock()
	// No path
	return Route{[]Node{}, -1}
}

// CalculatePath Finds closest nodes to start and end
func (graph *Graph) CalculatePath(startCoords Coordinate, endCoords Coordinate) Route {
	nodeStart := graph.FindNode(startCoords)
	nodeEnd := graph.FindNode(endCoords)
	fmt.Println("st, end,", nodeStart, nodeEnd)
	pathFound := graph.FindPath(nodeStart, nodeEnd)
	return pathFound
}
