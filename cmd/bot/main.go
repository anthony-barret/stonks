package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/anthony-barret/stonks/internal/discord"
	"github.com/anthony-barret/stonks/internal/providers/alphavantage"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	avAPIKey := os.Getenv("ALPHAVANTAGE_TOKEN")
	stockAPI := alphavantage.NewClient(avAPIKey)
	discordClient, _ := discord.NewClient(os.Getenv("DISCORD_APP_ID"), os.Getenv("DISCORD_GUILD_ID"), os.Getenv("DISCORD_TOKEN"), stockAPI)
	discordClient.Run()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	discordClient.Stop()
}
