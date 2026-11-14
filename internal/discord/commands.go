package discord

import (
	"sort"

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
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "symbol",
				Description: "Stock ticker symbol",
				Required:    true,
			},
		},
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

func (d *DiscordClient) commandDailyQuotation(i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	symbol := options[0].StringValue()
	timeSeriesDaily, err := d.stockAPIClient.GetTimeSeriesDailyForSymbol(symbol)
	if err != nil {
		d.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to fetch daily quotation for your ticker.",
			},
		})
		return
	}
	// maps do not preserve order, so reorder to get the latest keys
	keys := make([]string, 0, len(timeSeriesDaily.TimeSeries))
	for k := range timeSeriesDaily.TimeSeries {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	var embedFields []*discordgo.MessageEmbedField
	for j, day := range keys {
		if j > 9 {
			break
		}
		embedFields = append(embedFields, &discordgo.MessageEmbedField{
			Name:   day,
			Value:  "Open: " + timeSeriesDaily.TimeSeries[day].Open + " High: " + timeSeriesDaily.TimeSeries[day].High + " Low: " + timeSeriesDaily.TimeSeries[day].Low + " Close: " + timeSeriesDaily.TimeSeries[day].Close + " Volume: " + timeSeriesDaily.TimeSeries[day].Volume,
			Inline: false,
		})
	}
	d.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "Quotation for " + symbol,
				Description: "Last 10 days",
				Fields:      embedFields,
			}},
		},
	})
}
