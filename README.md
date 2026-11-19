# Stonks :chart_with_upwards_trend:

[![CI](https://github.com/anthony-barret/stonks/actions/workflows/ci.yml/badge.svg)](https://github.com/anthony-barret/stonks/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/anthony-barret/stonks/graph/badge.svg)](https://codecov.io/gh/anthony-barret/stonks)

**Stonks** is a simple Discord bot designed to track the stock market directly from your Discord server.

## Installation

### Configure environment variables

Create a `.env` file in the project root using `.env.example` as a template, then update it with your API credentials.

```sh
cp .env.example .env
```

### Start the bot

Run the following command from the project root:

```sh
go run ./cmd/bot
```

## Technologies

- [Go](https://go.dev/)
- [DiscordGo](https://github.com/bwmarrin/discordgo)
- [Alpha Vantage](https://www.alphavantage.co/)
