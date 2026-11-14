package discord

import (
	"testing"

	"github.com/anthony-barret/stonks/internal/providers/alphavantage"
)

func TestNewClient(t *testing.T) {
	applicationID := "appID"
	guildID := "guildID"
	token := "token"
	client, err := NewClient(applicationID, guildID, token, alphavantage.NewClient("test-api-key"))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if client.applicationID != applicationID {
		t.Errorf("expected applicationID %q, got %q", applicationID, client.applicationID)
	}
	if client.guildID != guildID {
		t.Errorf("expected guildID %q, got %q", guildID, client.guildID)
	}
	if client.session.Token != "Bot "+token {
		t.Errorf("expected session.Token %q, got %q", "Bot "+token, client.session.Token)
	}
	if _, ok := interface{}(client.stockAPIClient).(*alphavantage.Client); !ok {
		t.Fatal("stockAPIClient has an incorrect type")
	}
}
