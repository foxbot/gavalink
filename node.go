package gavalink

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// NodeConfig configures a Lavalink Node
type NodeConfig struct {
	// REST is the host where Lavalink's REST server runs
	//
	// This value is expected without a trailing slash, e.g. like
	// `http://localhost:2333`
	REST string
	// WebSocket is the host where Lavalink's WebSocket server runs
	//
	// This value is expected without a trailing slash, e.g. like
	// `http://localhost:8012`
	WebSocket string
	// Password is the expected Authorization header for the Node
	Password string
}

// Node wraps a Lavalink Node
type Node struct {
	config  NodeConfig
	load    float32
	manager *Lavalink
	wsConn  *websocket.Conn
}

func (node *Node) open() error {
	header := http.Header{}
	header.Set("Authorization", node.config.Password)
	header.Set("Num-Shards", node.manager.shards)
	header.Set("User-Id", node.manager.userID)

	ws, resp, err := websocket.DefaultDialer.Dial(node.config.WebSocket, header)
	if err != nil {
		return err
	}
	vstr := resp.Header.Get("Lavalink-Major-Version")
	v, err := strconv.Atoi(vstr)
	if err != nil {
		return err
	}
	if v < 3 {
		return errInvalidVersion
	}

	node.wsConn = ws
	go node.listen()

	Log.Println("node", node.config.WebSocket, "opened")

	return nil
}

func (node *Node) stop() {
	// someone already stopped this
	if node.wsConn == nil {
		return
	}
	_ = node.wsConn.Close()
}

func (node *Node) listen() {
	for {
		msgType, msg, err := node.wsConn.ReadMessage()
		if err != nil {
			Log.Println(err)
			// try to reconnect
			oerr := node.open()
			if oerr != nil {
				Log.Println("node", node.config.WebSocket, "failed and could not reconnect, destroying.", err, oerr)
				node.manager.removeNode(node)
				return
			}
			Log.Println("node", node.config.WebSocket, "reconnected")
			return
		}
		err = node.onEvent(msgType, msg)
		// TODO: better error handling?

		if err != nil {
			Log.Println(err)
		}
	}
}

func (node *Node) onEvent(msgType int, msg []byte) error {
	if msgType != websocket.TextMessage {
		return errUnknownPayload
	}

	m := message{}
	err := json.Unmarshal(msg, &m)
	if err != nil {
		return err
	}

	switch m.Op {
	case opPlayerUpdate:
		player, err := node.manager.GetPlayer(m.GuildID)
		if err != nil {
			return err
		}
		player.time = m.State.Time
		player.position = m.State.Position
	case opEvent:
		player, err := node.manager.GetPlayer(m.GuildID)
		if err != nil {
			return err
		}

		switch m.Type {
		case eventTrackEnd:
			player.track = ""
			err = player.handler.OnTrackEnd(player, m.Track, m.Reason)
		case eventTrackException:
			err = player.handler.OnTrackException(player, m.Track, m.Reason)
		case eventTrackStuck:
			err = player.handler.OnTrackStuck(player, m.Track, m.ThresholdMs)
		}

		return err
	case opStats:
		node.load = m.StatCPU.Load
	default:
		return errUnknownPayload
	}

	return nil
}

// CreatePlayer creates an audio player on this node
func (node *Node) CreatePlayer(guildID string, sessionID string, event VoiceServerUpdate, handler EventHandler) (*Player, error) {
	msg := message{
		Op:        opVoiceUpdate,
		GuildID:   guildID,
		SessionID: sessionID,
		Event:     &event,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	err = node.wsConn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return nil, err
	}
	player := &Player{
		guildID: guildID,
		manager: node.manager,
		node:    node,
		handler: handler,
		vol:     100,
	}
	node.manager.players[guildID] = player
	return player, nil
}

// LoadTracks queries lavalink to return a Tracks object
//
// query should be a valid Lavaplayer query, including but not limited to:
// - A direct media URI
// - A direct Youtube /watch URI
// - A search query, prefixed with ytsearch: or scsearch:
//
// See the Lavaplayer Source Code for all valid options.
func (node *Node) LoadTracks(query string) (*Tracks, error) {
	url := fmt.Sprintf("%s/loadtracks?identifier=%s", node.config.REST, query)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", node.config.Password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	tracks := new(Tracks)
	err = json.Unmarshal(data, &tracks)
	if err != nil {
		return nil, err
	}
	return tracks, nil
}
