package sqlite

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
	"vpn-tg-bot/internal/storage"
)

func TestAddUsers(t *testing.T) {
	db, err := New(path)
	if err != nil {
		t.Fatal("can't create db instance:", err)
	}

	time.Sleep(time.Millisecond * 200)

	n := 10
	for i := 1; i <= n; i++ {

		wantUser := &storage.User{
			TelegramID:   storage.TelegramID(i),
			TelegramName: "test" + strconv.Itoa(i),
		}

		// Add User
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		id, err := db.SaveUser(ctx, wantUser)
		cancel()
		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("added User with id %d\n", id)
		// time.Sleep(time.Microsecond * 100)

	}
}

func TestAddGetRemoveUsers(t *testing.T) {

	db, err := New(path)
	if err != nil {
		t.Fatal("can't create db instance:", err)
	}

	time.Sleep(time.Millisecond * 200)

	wantUser := &storage.User{
		TelegramID:   1,
		TelegramName: "test",
	}

	// Add User
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	id, err := db.SaveUser(ctx, wantUser)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("added User with id %d\n", id)
	time.Sleep(time.Millisecond * 100)

	// Get User
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	user, err := db.GetUserByID(ctx, id)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	if wantUser.TelegramID != user.TelegramID {
		t.Fatalf("want %d, got %d", wantUser.TelegramID, user.TelegramID)
	}
	if wantUser.TelegramName != user.TelegramName {
		t.Fatalf("want %s, got %s", wantUser.TelegramName, user.TelegramName)
	}

	fmt.Printf("User: %+v\n", user)
	time.Sleep(time.Millisecond * 100)

	// Remove User
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	err = db.RemoveUserByID(ctx, id)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	// Get User after removal
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	user, err = db.GetUserByID(ctx, id)
	cancel()
	if err == nil {
		t.Fatalf("user still exists: %+v", user)
	} else if err != storage.ErrNoSuchUser {
		t.Fatal(err)
	}

	fmt.Printf("removed User with id %d\n", id)
}

func TestGetUsers(t *testing.T) {
	db, err := New(path)
	if err != nil {
		t.Fatal("can't create db instance:", err)
	}

	time.Sleep(time.Millisecond * 200)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)

	args := &storage.QueryArgs{
		Offset: 0,
		Limit:  5,
	}
	users, err := db.GetUsers(ctx, args)
	if err != nil {
		t.Fatalf("can't get all users: %v", err)
	}
	cancel()
	t.Logf("got %d users", len(*users))
}
