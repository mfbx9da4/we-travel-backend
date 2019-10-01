package main

import (
	"container/heap"
	"fmt"
	"math"
	"strconv"

	//"math"
	"sync"
	//"github.com/cheekybits/genny/generic"
)

// Node a single node that composes the tree
type Node struct {
	Value Coordinate
}

func CreateNode(coords Coordinate) Node {
	var hash = strconv.FormatFloat(coords[0], 'f', -1, 64)
	hash += strconv.FormatFloat(coords[1], 'f', -1, 64)
	return Node{coords}
}

func (n *Node) String() string {
	return fmt.Sprintf("%v", n.Value)
}

type Graph struct {
	nodes []*Node
	edges map[Node][]*Node
	lock  sync.RWMutex
}

// AddNode adds a node to the graph
func (graph *Graph) AddNode(n *Node) {
	graph.lock.Lock()
	graph.nodes = append(graph.nodes, n)
	graph.lock.Unlock()
}

// AddEdge adds an edge to the graph
func (graph *Graph) AddEdge(n1, n2 *Node) {
	graph.lock.Lock()
	if graph.edges == nil {
		graph.edges = make(map[Node][]*Node)
	}
	graph.edges[*n1] = append(graph.edges[*n1], n2)
	graph.edges[*n2] = append(graph.edges[*n2], n1)
	graph.lock.Unlock()
}

// Print graph
func (graph *Graph) String() {
	graph.lock.RLock()
	s := ""
	for i := 0; i < len(graph.nodes); i++ {
		s += graph.nodes[i].String() + " -> "
		near := graph.edges[*graph.nodes[i]]
		for j := 0; j < len(near); j++ {
			s += near[j].String() + " "
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

func (graph *Graph) FindNode(coords Coordinate) int {
	nodes := graph.nodes
	//min distance
	var foundNodeIndex int
	var minDistance float64
	// Probably there is a better algo for this, just doing the brute force sorry :(
	for i := 0; i < len(nodes); i++ {
		//calc distance
		value := nodes[i].Value
		dx := coords[0] - value[0]
		dy := coords[1] - value[1]
		distance := math.Sqrt(dx*dx + dy*dy)
		if i == 0 || distance < minDistance {
			foundNodeIndex = i
			minDistance = distance
		}
	}
	return foundNodeIndex
}

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

		for i := 0; i < len(children); i++ {
			child := children[i]
			dx := (node.Value[0] - child.Value[0])
			dy := (node.Value[1] - child.Value[1])
			remaingDx := (dest.Value[0] - child.Value[0])
			remainingDy := (dest.Value[1] - child.Value[1])
			elapsed := math.Sqrt(dx*dx+dy*dy) + cur.Distance
			remaining := math.Sqrt(remaingDx*remaingDx + remainingDy*remainingDy)

			if *child == *dest {
				fmt.Println("Found Dest with distance", pqitem.Priority)
				path := append(cur.Path, *child)
				return Route{path, elapsed}
			}

			if !visited[child] {
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

func (graph *Graph) CalculatePath(startCoords Coordinate, endCoords Coordinate) Route {
	nodeStart := graph.FindNode(startCoords)
	nodeEnd := graph.FindNode(endCoords)
	pathFound := graph.FindPath(graph.nodes[nodeStart], graph.nodes[nodeEnd])
	return pathFound
}
