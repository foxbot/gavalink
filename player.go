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
	paused   bool
	vol      int
	track    string
	manager  *Lavalink
	node     *Node
	handler  EventHandler
}

// GuildID returns this player's Guild ID
func (player *Player) GuildID() string {
	return player.guildID
}

// Play will play the given track completely
func (player *Player) Play(track string) error {
	return player.PlayAt(track, 0, 0)
}

// PlayAt will play the given track at the specified start and end times
//
// Setting a time to 0 will omit it.
func (player *Player) PlayAt(track string, startTime int, endTime int) error {
	player.paused = false
	player.track = track

	start := strconv.Itoa(startTime)
	end := strconv.Itoa(endTime)

	msg := message{
		Op:        opPlay,
		GuildID:   player.guildID,
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

// Track returns the player's current track
func (player *Player) Track() string {
	return player.track
}

// Stop will stop the currently playing track
func (player *Player) Stop() error {
	player.track = ""
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
	player.paused = pause

	msg := message{
		Op:      opPause,
		GuildID: player.guildID,
		Pause:   &pause,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = player.node.wsConn.WriteMessage(websocket.TextMessage, data)
	return err
}

// Paused returns whether or not the player is currently paused
func (player *Player) Paused() bool {
	return player.paused
}

// Seek will seek the player to the speicifed position, in millis
func (player *Player) Seek(position int) error {
	msg := message{
		Op:       opSeek,
		GuildID:  player.guildID,
		Position: &position,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = player.node.wsConn.WriteMessage(websocket.TextMessage, data)
	return err
}

// Position returns the player's position, as reported by Lavalink
func (player *Player) Position() int {
	return player.position
}

// Volume will set the player's volume to the specified value
//
// volume must be within [0, 1000]
func (player *Player) Volume(volume int) error {
	if volume < 0 || volume > 1000 {
		return errVolumeOutOfRange
	}

	player.vol = volume

	msg := message{
		Op:      opVolume,
		GuildID: player.guildID,
		Volume:  &volume,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = player.node.wsConn.WriteMessage(websocket.TextMessage, data)
	return err
}

// GetVolume gets the player's volume level
func (player *Player) GetVolume() int {
	return player.vol
}

// Forward will forward a new VOICE_SERVER_UPDATE to a Lavalink node for
// this player.
//
// This should always be used if a VOICE_SERVER_UPDATE is received for
// a guild which already has a player.
//
// To move a player to a new Node, first player.Destroy() it, and then
// create a new player on the new node.
func (player *Player) Forward(sessionID string, event VoiceServerUpdate) error {
	msg := message{
		Op:        opVoiceUpdate,
		GuildID:   player.guildID,
		SessionID: sessionID,
		Event:     &event,
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
	delete(player.manager.players, player.guildID)
	return nil
}
