package sqlite

import (
	"context"
	"fmt"
	"testing"
	"time"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/structconv"
)

func TestSaveGetDeleteAuthorization(t *testing.T) {
	db, err := New(path)
	if err != nil {
		t.Fatal("can't create db instance:", err)
	}

	time.Sleep(time.Millisecond * 200)

	testAuth := &storage.VPNServerAuthorization{
		ServerID:  100,
		ExpiredAt: time.Now().Add(time.Hour),
		Username:  "test",
		Password:  "test",
		Token:     "test",
		Meta:      "test",
	}

	// Save Authorization
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	id, err := db.SaveServerAuth(ctx, testAuth)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("added Authorization with id %d\n", id)
	time.Sleep(time.Millisecond * 100)

	// Test Update
	testAuth.ExpiredAt = testAuth.ExpiredAt.Add(time.Hour)
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	id, err = db.SaveServerAuth(ctx, testAuth)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("updated Authorization with id %d\n", id)
	time.Sleep(time.Millisecond * 100)

	// Get Authorization
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	auth, err := db.GetServerAuthByServerID(ctx, id)
	cancel()
	if err != nil {
		t.Fatal(err)
	}
	if err = structconv.CompareStructs(auth, testAuth, []string{"ServerID", "Username", "Password", "Token", "ExpiredAt", "Meta"}); err != nil {
		t.Fatal(err)
	}
	t.Log(auth)

	// Delete Authorization
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	err = db.RemoveServerAuthByServerID(ctx, id)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	// Get Server after removal
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	auth, err = db.GetServerAuthByServerID(ctx, id)
	cancel()
	if err == nil {
		t.Fatalf("server still exists: %+v", auth)
	} else if err != storage.ErrNoSuchServerAuth {
		t.Fatal(err)
	}

	fmt.Printf("removed Authorization with id %d\n", id)
}
