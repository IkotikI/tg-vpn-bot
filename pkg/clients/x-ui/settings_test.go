package x_ui

import (
	"context"
	"testing"
	"time"
	"vpn-tg-bot/pkg/debug"
)

func TestGetAllSettings(t *testing.T) {
	xui := New(TokenKey_3x_ui, server, makeAuthStore())

	// Add Inbound
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	settings, err := xui.GetAllSettings(ctx)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("settings:\n%s", debug.JSON(settings))
}
