package x_ui

import (
	"context"
	"fmt"
	"testing"
	"time"
	"vpn-tg-bot/pkg/clients/x-ui/model"
	"vpn-tg-bot/pkg/debug"
)

var test_inboundID = 1

var test_email = "test@example.com"
var test_uuid = "f00ec978-7f3a-49b0-aff6-fd3c88fa6186"

// var test_email = "2test@example.com"
// var test_uuid = "200ec978-7f3a-49b0-aff6-fd3c88fa6186"

func TestAddUpdateRemoveClient(t *testing.T) {
	TestAddClient(t)
	TestUpdateClient(t)
	TestDeleteClient(t)
}

func TestAddClient(t *testing.T) {
	xui := New(TokenKey_3x_ui, server, makeAuthStore())

	client := &model.Client{
		ID:         test_uuid,
		Email:      test_email,
		ExpiryTime: time.Now().Add(time.Hour * 24).Unix(),
	}

	bytes_test, err := xui.prepareClientPayload(test_inboundID, client)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("payload:\n%+v\n", string(bytes_test))

	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	err = xui.AddClient(ctx, test_inboundID, client)
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	time.Sleep(time.Millisecond * 100)
}

func TestUpdateClient(t *testing.T) {
	xui := New(TokenKey_3x_ui, server, makeAuthStore())

	client := &model.Client{
		ID:         test_uuid,
		Email:      test_email,
		ExpiryTime: time.Now().Add(time.Hour * 24).Unix(),
	}

	bytes_test, err := xui.prepareClientPayload(test_inboundID, client)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("payload:\n%+v\n", string(bytes_test))

	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	err = xui.UpdateClient(ctx, test_inboundID, client)
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	time.Sleep(time.Millisecond * 100)
}

func TestDeleteClient(t *testing.T) {
	xui := New(TokenKey_3x_ui, server, makeAuthStore())

	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	clientID, _ := ParseClientID(test_uuid)
	err := xui.DeleteClient(ctx, test_inboundID, clientID)
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	fmt.Println("client deleted successfully")
}

func TestGetClientTraffic(t *testing.T) {
	xui := New(TokenKey_3x_ui, server, makeAuthStore())

	clientID, _ := ParseClientID(test_uuid)
	// Add test user
	testClient := &model.Client{
		ID:         test_uuid,
		Email:      test_email,
		ExpiryTime: time.Now().Add(time.Hour * 24).Unix(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	if err := xui.AddClient(ctx, test_inboundID, testClient); err != nil {
		t.Log(err)
	}
	cancel()
	time.Sleep(time.Millisecond * 100)

	var clientByEmail *model.ClientTraffic
	var clientByID *model.ClientTraffic
	var err error
	t.Run("GetClientByEmail", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		clientByEmail, err = xui.GetClientClientTrafficsByEmail(ctx, test_email)
		if err != nil {
			t.Fatal(err)
		}
		cancel()
		t.Logf("client traffic:\n%s\n", debug.JSON(clientByEmail))
	})
	time.Sleep(time.Millisecond * 100)

	t.Run("GetClientByID", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		clients, err := xui.GetClientTrafficByID(ctx, clientID)
		if err != nil {
			t.Fatal(err)
		}
		cancel()
		t.Logf("clients traffic:\n%s\n", debug.JSON(clients))
		for _, client := range *clients {
			if client.Email == test_email {
				clientByID = &client
				break
			}
		}
	})
	time.Sleep(time.Millisecond * 100)

	ctx, cancel = context.WithTimeout(context.Background(), 400*time.Millisecond)
	if err := xui.DeleteClient(ctx, test_inboundID, clientID); err != nil {
		t.Log(err)
	}
	cancel()
	time.Sleep(time.Millisecond * 100)

	if *clientByEmail != *clientByID {
		t.Fatalf("Client's traffic is not equal between methods: want: %+v, got: %+v", *clientByID, *clientByEmail)
	}

}
