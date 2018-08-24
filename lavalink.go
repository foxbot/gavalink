package gavalink

import (
	"errors"
	"log"
	"os"
	"sort"
)

// Log sets the log.Logger gavalink will write to
var Log *log.Logger

func init() {
	Log = log.New(os.Stdout, "(gavalink) ", 0)
}

// Lavalink manages a connection to Lavalink Nodes
type Lavalink struct {
	shards string
	userID string

	nodes   []Node
	players map[string]*Player
}

var (
	errNoNodes          = errors.New("No nodes present")
	errNodeNotFound     = errors.New("Couldn't find that node")
	errPlayerNotFound   = errors.New("Couldn't find a player for that guild")
	errVolumeOutOfRange = errors.New("Volume is out of range, must be within [0, 1000]")
	errInvalidVersion   = errors.New("This library requires Lavalink >= 3")
	errUnknownPayload   = errors.New("Lavalink sent an unknown payload")
	errNilHandler       = errors.New("You must provide an event handler. Use gavalink.DummyEventHandler if you wish to ignore events")
)

// NewLavalink creates a new Lavalink manager
func NewLavalink(shards string, userID string) *Lavalink {
	return &Lavalink{
		shards: shards,
		userID: userID,
		/*		nodes:   make([]Node, 1),*/
		players: make(map[string]*Player),
	}
}

// AddNodes adds a node to the Lavalink manager
func (lavalink *Lavalink) AddNodes(nodeConfigs ...NodeConfig) error {
	nodes := make([]Node, len(nodeConfigs))
	for i, c := range nodeConfigs {
		n := Node{
			config:  c,
			manager: lavalink,
		}
		err := n.open()
		if err != nil {
			return err
		}
		nodes[i] = n
	}
	lavalink.nodes = append(lavalink.nodes, nodes...)
	return nil
}

// RemoveNode removes a node from the manager
func (lavalink *Lavalink) removeNode(node *Node) error {
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

	sort.SliceStable(lavalink.nodes, func(i, j int) bool {
		return lavalink.nodes[i].load < lavalink.nodes[j].load
	})

	return &lavalink.nodes[0], nil
}

// GetPlayer gets a player for a guild
func (lavalink *Lavalink) GetPlayer(guild string) (*Player, error) {
	p, ok := lavalink.players[guild]
	if !ok {
		return nil, errPlayerNotFound
	}
	return p, nil
}
