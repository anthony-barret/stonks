package discord

import (
	"fmt"

	"github.com/anthony-barret/stonks/internal/providers/alphavantage"
	"github.com/bwmarrin/discordgo"
)

type DiscordClient struct {
	session        *discordgo.Session
	applicationID  string
	guildID        string
	stockAPIClient *alphavantage.Client
}

func NewClient(appID string, guildID string, token string, stockAPIClient *alphavantage.Client) (*DiscordClient, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("create discord session: %w", err)
	}

	client := &DiscordClient{session: session, applicationID: appID, guildID: guildID, stockAPIClient: stockAPIClient}
	client.session.AddHandler(client.handleInteraction)

	return client, nil
}

func (d *DiscordClient) Run() error {
	if err := d.session.Open(); err != nil {
		return fmt.Errorf("open discord session: %w", err)
	}
	_, err := d.session.ApplicationCommandBulkOverwrite(d.applicationID, d.guildID, commands)
	if err != nil {
		d.session.Close()
		return fmt.Errorf("register discord commands: %w", err)
	}
	return nil
}

func (d *DiscordClient) Stop() error {
	if err := d.session.Close(); err != nil {
		return fmt.Errorf("close discord session: %w", err)
	}
	return nil
}
