package discord

import (
	"fmt"
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

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		data := i.ApplicationCommandData()
		switch data.Name {
		case "market-status":
			commandMarketStatus(s, i)
		}
	})
	return &DiscordClient{session: session, applicationID: appID, guildID: guildID, stockAPIClient: stockAPIClient}
}

func (d *DiscordClient) Run() error {
	if err := d.session.Open(); err != nil {
		return fmt.Errorf("open Discord session: %w", err)
	}
	_, err := d.session.ApplicationCommandBulkOverwrite(d.applicationID, d.guildID, commands)
	if err != nil {
		d.session.Close()
		return fmt.Errorf("register Discord commands: %w", err)
	}
	return nil
}

func (d *DiscordClient) Stop() error {
	if err := d.session.Close(); err != nil {
		return fmt.Errorf("close Discord session: %w", err)
	}
	return nil
}
