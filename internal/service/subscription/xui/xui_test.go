package xui_service

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/internal/storage/sqlite"
	x_ui "vpn-tg-bot/pkg/clients/x-ui"
)

const basePath = "../../../../"
const internalPath = "../../../"

var testStoragePath = internalPath + "storage/sqlite/test_data/db.db"
var test_UserID = storage.UserID(1)
var test_ServerID = storage.ServerID(1)

func makeStore() *sqlite.SQLStorage {
	db, err := sqlite.New(testStoragePath)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Millisecond * 200)
	return db
}

func TestUpdateSubscription(t *testing.T) {

	db := makeStore()

	subscription := &storage.Subscription{
		UserID:                test_UserID,
		ServerID:              test_ServerID,
		SubscriptionStatus:    storage.SubscriptionStatusActive,
		SubscriptionExpiredAt: time.Now().Add(time.Hour),
	}

	s := NewXUIService(x_ui.TokenKey_3x_ui, db, db)

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	err := s.UpdateSubscription(ctx, subscription)
	if err != nil {
		t.Fatal(err)
	}
	cancel()

}

func TestGetSub(t *testing.T) {
	db := makeStore()

	s := NewXUIService(x_ui.TokenKey_3x_ui, db, db)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	link, err := s.SubscriptionLink(ctx, test_ServerID, test_UserID)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("link: %s", link)

}

func TestGetSubDirect(t *testing.T) {
	// db := makeStore()

	// s := NewXUIService(x_ui.TokenKey_3x_ui, db, db)
	// ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	// settigs, err := s.xuiClientInstance(ctx, test_ServerID)
	// cancel()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	subId := "n2b9ubaioe06cak8"

	// ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	resp, err := http.Get("http://localhost:2096/sub/" + subId)
	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		t.Fatal(err)
	}

	link := string(bytes)
	t.Log(link)

}
