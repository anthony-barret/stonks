package discord

import (
	"log"

	"github.com/anthony-barret/stonks/internal/providers/alphavantage"
	"github.com/bwmarrin/discordgo"
)

type DiscordClient struct {
	session        *discordgo.Session
	applicationID  string
	guildID        string
	stockAPIClient *alphavantage.Client
}

func NewClient(appID string, guildID string, token string, stockAPIClient *alphavantage.Client) *DiscordClient {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error during the creation of the discord bot", err)
	}
	session.ApplicationCommandBulkOverwrite(appID, guildID, []*discordgo.ApplicationCommand{
		{
			Name:        "market-status",
			Description: "Check the market status",
		},
	})

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		data := i.ApplicationCommandData()
		switch data.Name {
		case "market-status":
			commandMarketStatus(s, i)
		}
	})
	return &DiscordClient{session: session, applicationID: appID, guildID: guildID, stockAPIClient: stockAPIClient}
}

func (d *DiscordClient) Run() {
	err := d.session.Open()
	if err != nil {
		log.Fatal("Error when opening the discord session", err)
	}
}

func (d *DiscordClient) Stop() {
	err := d.session.Close()
	if err != nil {
		log.Fatal("Error when closing the discord session", err)
	}
}
