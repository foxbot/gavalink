package gavalink

import (
	"errors"
)

// Lavalink manages a connection to Lavalink Nodes
type Lavalink struct {
	// Shards is the total number of shards the bot is running
	Shards int
	// UserID is the Discord User ID of the bot
	UserID int

	nodes []Node
}

var (
	errNoNodes        = errors.New("No nodes present")
	errNodeNotFound   = errors.New("Couldn't find that node")
	errInvalidVersion = errors.New("This library requires Lavalink >= 3")
	errUnknownPayload = errors.New("Lavalink sent an unknown payload")
)

// NewLavalink creates a new Lavalink manager
func NewLavalink() *Lavalink {
	return &Lavalink{}
}

// AddNodes adds a node to the Lavalink manager
func (lavalink *Lavalink) AddNodes(nodeConfigs ...NodeConfig) {
	nodes := make([]Node, len(nodeConfigs))
	for i, c := range nodeConfigs {
		n := Node{
			config: c,
			shards: lavalink.Shards,
			userID: lavalink.UserID,
		}
		nodes[i] = n
	}
	lavalink.nodes = append(lavalink.nodes, nodes...)
}

// RemoveNode removes a node from the manager
func (lavalink *Lavalink) RemoveNode(node *Node) error {
	idx := -1
	for i, n := range lavalink.nodes {
		if n == *node {
			idx = i
			break
		}
	}
	if idx == -1 {
		return errNodeNotFound
	}

	node.stop()

	// temp var for easier reading
	n := lavalink.nodes
	z := len(n) - 1

	n[idx] = n[z] // swap idx with last
	n = n[:z]

	lavalink.nodes = n
	return nil
}

// BestNode returns the Node with the lowest latency
func (lavalink *Lavalink) BestNode() (*Node, error) {
	if len(lavalink.nodes) < 1 {
		return nil, errNoNodes
	}
	// TODO: lookup latency
	return &lavalink.nodes[0], nil
}
