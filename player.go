package gavalink

import (
	"encoding/json"
	"strconv"

	"github.com/gorilla/websocket"
)

// Player is a Lavalink player
type Player struct {
	guildID  string
	time     int
	position int
	manager  *Lavalink
	node     *Node
	handler  EventHandler
}

// Play will play the given track completely
func (player *Player) Play(track string) error {
	return player.PlayAt(track, 0, 0)
}

// PlayAt will play the given track at the specified start and end times
//
// Setting a time to 0 will omit it.
func (player *Player) PlayAt(track string, startTime int, endTime int) error {
	start := strconv.Itoa(startTime)
	end := strconv.Itoa(endTime)

	msg := message{
		Op:        opPlay,
		Track:     track,
		StartTime: start,
		EndTime:   end,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = player.node.wsConn.WriteMessage(websocket.TextMessage, data)
	return err
}

// Stop will stop the currently playing track
func (player *Player) Stop() error {
	msg := message{
		Op:      opStop,
		GuildID: player.guildID,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = player.node.wsConn.WriteMessage(websocket.TextMessage, data)
	return err
}

// Pause will pause or resume the player, depending on the pause parameter
func (player *Player) Pause(pause bool) error {
	msg := message{
		Op:      opPause,
		GuildID: player.guildID,
		Pause:   pause,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = player.node.wsConn.WriteMessage(websocket.TextMessage, data)
	return err
}

// Seek will seek the player to the speicifed position, in millis
func (player *Player) Seek(position int) error {
	msg := message{
		Op:       opSeek,
		GuildID:  player.guildID,
		Position: position,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = player.node.wsConn.WriteMessage(websocket.TextMessage, data)
	return err
}

// Volume will set the player's volume to the specified value
//
// volume must be within [0, 1000]
func (player *Player) Volume(volume int) error {
	if volume < 0 || volume > 1000 {
		return errVolumeOutOfRange
	}

	msg := message{
		Op:      opVolume,
		GuildID: player.guildID,
		Volume:  volume,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = player.node.wsConn.WriteMessage(websocket.TextMessage, data)
	return err
}

// Destroy will destroy this player
func (player *Player) Destroy() error {
	msg := message{
		Op:      opDestroy,
		GuildID: player.guildID,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = player.node.wsConn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}
	player.manager.players[player.guildID] = nil
	return nil
}
