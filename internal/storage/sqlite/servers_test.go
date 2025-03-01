package sqlite

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strconv"
	"testing"
	"time"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/structconv"
)

func TestAddGetRemoveServers(t *testing.T) {

	db, err := New(path)
	if err != nil {
		t.Fatal("can't create db instance:", err)
	}

	time.Sleep(time.Millisecond * 200)

	testServer := &storage.VPNServer{
		CountryID: 1,
		Name:      "test",
		Protocol:  "test",
		Host:      "test",
		Port:      1,
		Username:  "test",
		Password:  "test",
	}

	// Add Server
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	id, err := db.SaveServer(ctx, testServer)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("added Server with id %d\n", id)
	time.Sleep(time.Millisecond * 100)

	// Get Server
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	server, err := db.GetServerByID(ctx, id)
	cancel()
	if err != nil {
		t.Fatal(err)
	}
	if err = structconv.CompareStructs(server, testServer, []string{"CountryID", "Name", "Protocol", "IPaddress", "Port", "Login", "Password"}); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Server: %+v\n", server)
	time.Sleep(time.Millisecond * 100)

	// Remove Server
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	err = db.RemoveServerByID(ctx, id)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	// Get Server after removal
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	server, err = db.GetServerByID(ctx, id)
	cancel()
	if err == nil {
		t.Fatalf("server still exists: %+v", server)
	} else if err != storage.ErrNoSuchServer {
		t.Fatal(err)
	}

	fmt.Printf("removed Server with id %d\n", id)
}

func TestAddServers(t *testing.T) {
	db, err := New(path)
	if err != nil {
		t.Fatal("can't create db instance:", err)
	}

	time.Sleep(time.Millisecond * 200)

	n := 10
	for i := 0; i < n; i++ {

		testServer := &storage.VPNServer{
			CountryID: storage.CountryID(1 + rand.Int64N(250)),
			Name:      "test" + strconv.Itoa(i),
			Protocol:  "test" + strconv.Itoa(i),
			Host:      "test" + strconv.Itoa(i),
			Port:      i,
			Username:  "test" + strconv.Itoa(i),
			Password:  "test" + strconv.Itoa(i),
		}

		// Add Server
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		id, err := db.SaveServer(ctx, testServer)
		cancel()
		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("added Server with id %d\n", id)
		// time.Sleep(time.Microsecond * 100)

	}
}

func TestUpdateServer(t *testing.T) {
	s, err := New(path)
	if err != nil {
		t.Fatal(err)
	}

	args := &storage.QueryArgs{
		Where: []storage.Where{
			{
				Column:   "country_id",
				Operator: storage.OpMore,
				Value:    250,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	servers, err := s.GetServers(ctx, args)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	for _, server := range *servers {
		fmt.Printf("id: %d, country_id %d\n", server.ID, server.CountryID)
	}

	// TODO
}
