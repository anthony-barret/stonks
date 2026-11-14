package discord

import (
	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "market-status",
		Description: "Check the market status",
	},
	{
		Name:        "daily-quotation",
		Description: "Get the daily quotation for a symbol",
	},
}

func (d *DiscordClient) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "market-status":
		d.commandMarketStatus(i)

	case "daily-quotation":
		d.commandDailyQuotation(i)
	}
}

func (d *DiscordClient) commandMarketStatus(i *discordgo.InteractionCreate) {
	marketStatus, err := d.stockAPIClient.GetMarketStatus()
	if err != nil {
		d.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to fetch market status.",
			},
		})
		return
	}
	var embedFields []*discordgo.MessageEmbedField
	for _, market := range marketStatus.Markets {
		embedFields = append(embedFields, &discordgo.MessageEmbedField{
			Name:   market.MarketType + ": " + market.PrimaryExchanges,
			Value:  market.CurrentStatus,
			Inline: true,
		})
	}
	d.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "Market Status",
				Description: marketStatus.Endpoint,
				Fields:      embedFields,
			}},
		},
	})
}
