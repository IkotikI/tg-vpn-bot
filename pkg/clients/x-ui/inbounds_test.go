package x_ui

import (
	"context"
	"fmt"
	"testing"
	"time"
	"vpn-tg-bot/pkg/clients/x-ui/model"
)

var test_inboundID_ = 5

func TestLogin(t *testing.T) {
	xui := New(TokenKey_3x_ui, server, makeAuthStore())

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	xui.Auth(ctx)
}

func TestAddInbound(t *testing.T) {
	xui := New(TokenKey_3x_ui, server, makeAuthStore())

	// Add Inbound
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	var inboundToAdd model.Inbound = DefaultInbound
	// inboundToAdd.Id = test_inboundID_
	inbound, err := xui.AddInbound(ctx, &inboundToAdd)
	if err != nil {
		t.Fatalf("can't add inbound: %s", err)
	}
	fmt.Printf("added inbound: %+v\n", inbound)
	cancel()
	time.Sleep(time.Millisecond * 100)
}

func TestGetUpdateInbound(t *testing.T) {
	xui := New(TokenKey_3x_ui, server, makeAuthStore())

	// ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	// err = xui.Auth(ctx)
	// if err != nil {
	// 	t.Fatalf("can't auth: %s", err)
	// }
	// cancel()

	// Get Inbound
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	inbound, err := xui.GetInbound(ctx, test_inboundID_)
	if err != nil {
		t.Fatalf("can't get inbound: %s", err)
	}
	fmt.Printf("got inbound: %+v\n", inbound)
	cancel()
	time.Sleep(time.Millisecond * 100)

	// Update Inbound
	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	updInbound, err := xui.UpdateInbound(ctx, inbound)
	if err != nil {
		t.Fatalf("can't update inbound: %s", err)
	}
	fmt.Printf("updated inbound: %+v\n", updInbound)
	cancel()
	time.Sleep(time.Millisecond * 100)

}

// func TestUpdateInbound(t *testing.T) {
// 	xui := New(TokenKey_3x_ui, server, makeAuthStore())

// 	// Update Inbound
// 	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
// 	updInbound, err := xui.UpdateInbound(ctx, inbound)
// 	if err != nil {
// 		t.Fatalf("can't update inbound: %s", err)
// 	}
// 	fmt.Printf("updated inbound: %+v\n", updInbound)
// 	cancel()
// 	time.Sleep(time.Millisecond * 100)
// }

func TestDeleteInbound(t *testing.T) {
	xui := New(TokenKey_3x_ui, server, makeAuthStore())

	// Remove Inbound
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	err := xui.DeleteInbound(ctx, test_inboundID_)
	if err != nil {
		t.Fatalf("can't remove inbound: %s", err)
	}
	cancel()
	time.Sleep(time.Millisecond * 100)
}
