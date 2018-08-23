package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/foxbot/gavalink"

	// this requires the "feature/manual-voice-connection branch!!"
	"github.com/foxbot/discordgo"
)

var token string
var lavalink *gavalink.Lavalink
var player *gavalink.Player

func init() {
	flag.StringVar(&token, "token", "", "token=unprefixed token")
}

func main() {
	flag.Parse()

	if token == "" {
		panic("no token specified!")
	}
	token = "Bot " + token

	dg, err := discordgo.New(token)
	if err != nil {
		panic(err)
	}
	dg.SyncEvents = false

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)
	dg.AddHandler(voiceServerUpdate)

	err = dg.Open()
	if err != nil {
		panic(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("discordgo ready!")
	s.UpdateStatus(0, "gavalink")

	lavalink = gavalink.NewLavalink("1", event.User.ID)

	err := lavalink.AddNodes(gavalink.NodeConfig{
		REST:      "http://localhost:2333",
		WebSocket: "ws://localhost:2334",
		Password:  "youshallnotpass",
	})
	if err != nil {
		log.Println(err)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "~>>join" {
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			log.Println("fail find channel")
			return
		}

		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			log.Println("fail find guild")
			return
		}

		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				log.Println("trying to connect to channel")
				err = s.ChannelVoiceJoinManual(c.GuildID, vs.ChannelID, false, false)
				if err != nil {
					log.Println(err)
				} else {
					log.Println("channel voice join succeeded")
				}
			}
		}
	} else if m.Content == "~>>playtest" {
		node, err := lavalink.BestNode()
		if err != nil {
			log.Println(err)
			return
		}
		tracks, err := node.LoadTracks("https://www.youtube.com/watch?v=Wl2Q_MceIUc")
		if err != nil {
			log.Println(err)
			return
		}
		if tracks.Type != gavalink.TrackLoaded {
			log.Println("weird tracks type", tracks.Type)
		}
		track := tracks.Tracks[0].Data
		err = player.Play(track)
		if err != nil {
			log.Println(err)
		}
	}
}

func voiceServerUpdate(s *discordgo.Session, event *discordgo.VoiceServerUpdate) {
	log.Println("received VSU")
	node, err := lavalink.BestNode()
	if err != nil {
		log.Println(err)
		return
	}

	vsu := gavalink.VoiceServerUpdate{
		Endpoint: event.Endpoint,
		GuildID:  event.GuildID,
		Token:    event.Token,
	}
	handler := new(gavalink.DummyEventHandler)
	player, err = node.CreatePlayer(event.GuildID, s.State.SessionID, vsu, handler)
	if err != nil {
		log.Println(err)
		return
	}
}
