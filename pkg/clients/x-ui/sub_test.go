package x_ui

import (
	"context"
	"testing"
	"time"
)

var subID = "n2b9ubaioe06cak8"

func TestGetSubLink(t *testing.T) {
	xui := New(TokenKey_3x_ui, server, makeAuthStore())

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	link, err := xui.GetSubBySubID(ctx, subID)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("sub link:\n%s\n", link)
}
