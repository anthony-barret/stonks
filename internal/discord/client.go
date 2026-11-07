package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type DiscordClient struct {
	session *discordgo.Session
}

func NewClient(token string) *DiscordClient {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error during the creation of the discord bot", err)
	}
	return &DiscordClient{session: session}
}

func (d *DiscordClient) Run() {
	err := d.session.Open()
	if err != nil {
		log.Fatal("Error during the opening of the discord session", err)
	}
}
