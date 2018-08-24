package gavalink

const (
	// TrackLoaded is a Tracks Type for a succesful single track load
	TrackLoaded = "TRACK_LOADED"
	// PlaylistLoaded is a Tracks Type for a succseful playlist load
	PlaylistLoaded = "PLAYLIST_LOADED"
	// SearchResult is a Tracks Type for a search containing many tracks
	SearchResult = "SEARCH_RESULT"
	// NoMatches is a Tracks Type for a query yielding no matches
	NoMatches = "NO_MATCHES"
	// LoadFailed is a Tracks Type for an internal Lavalink error
	LoadFailed = "LOAD_FAILED"
)

// Tracks contains data for a Lavalink Tracks response
type Tracks struct {
	// Type contains the type of response
	//
	// This will be one of TrackLoaded, PlaylistLoaded, SearchResult,
	// NoMatches, or LoadFailed
	Type         string        `json:"loadType"`
	PlaylistInfo *PlaylistInfo `json:"playlistInfo"`
	Tracks       []Track       `json:"tracks"`
}

// PlaylistInfo contains information about a loaded playlist
type PlaylistInfo struct {
	// Name is the friendly of the playlist
	Name string `json:"name"`
	// SelectedTrack is the index of the track that loaded the playlist,
	// if one is present.
	SelectedTrack int `json:"selectedTrack"`
}

// Track contains information about a loaded track
type Track struct {
	// Data contains the base64 encoded Lavaplayer track
	Data string    `json:"track"`
	Info TrackInfo `json:"info"`
}

// TrackInfo contains more data about a loaded track
type TrackInfo struct {
	Identifier string `json:"identifier"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	URI        string `json:"uri"`
	Seekable   bool   `json:"isSeekable"`
	Stream     bool   `json:"isStream"`
	Length     int    `json:"length"`
	Position   int    `json:"position"`
}

const (
	opVoiceUpdate       = "voiceUpdate"
	opPlay              = "play"
	opStop              = "stop"
	opPause             = "pause"
	opSeek              = "seek"
	opVolume            = "volume"
	opDestroy           = "destroy"
	opPlayerUpdate      = "playerUpdate"
	opEvent             = "event"
	opStats             = "stats"
	eventTrackEnd       = "TrackEndEvent"
	eventTrackException = "TrackExceptionEvent"
	eventTrackStuck     = "TrackStuckEvent"
)

type message struct {
	Op          string             `json:"op"`
	GuildID     string             `json:"guildId,omitempty"`
	SessionID   string             `json:"sessionId,omitempty"`
	Event       *VoiceServerUpdate `json:"event,omitempty"`
	Track       string             `json:"track,omitempty"`
	StartTime   string             `json:"startTime,omitempty"`
	EndTime     string             `json:"endTime,omitempty"`
	Pause       *bool              `json:"pause,omitempty"`
	Position    *int               `json:"position,omitempty"`
	Volume      *int               `json:"volume,omitempty"`
	State       *state             `json:"state,omitempty"`
	Type        string             `json:"type,omitempty"`
	Reason      string             `json:"reason,omitempty"`
	Error       string             `json:"error,omitempty"`
	ThresholdMs int                `json:"thresholdMs,omitempty"`
	StatCPU     *statCPU           `json:"cpu,omitempty"`
	// TODO: stats
}

type state struct {
	Time     int `json:"time"`
	Position int `json:"position"`
}

type statCPU struct {
	Load float32 `json:"lavalinkLoad"`
}

// VoiceServerUpdate is a raw Discord VOICE_SERVER_UPDATE event
type VoiceServerUpdate struct {
	GuildID  string `json:"guild_id"`
	Endpoint string `json:"endpoint"`
	Token    string `json:"token"`
}
