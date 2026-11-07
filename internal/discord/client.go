package discord

import (
	"github.com/bwmarrin/discordgo"
)

type DiscordClient struct {
	session *discordgo.Session
}
