package discord

import (
	"os"

	"github.com/anthony-barret/stonks/internal/providers/alphavantage"
	"github.com/bwmarrin/discordgo"
)

func commandMarketStatus(s *discordgo.Session, i *discordgo.InteractionCreate) {
	api := alphavantage.NewClient(os.Getenv("ALPHAVANTAGE_TOKEN"))
	marketStatus, err := api.GetMarketStatus()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
