package sqlite

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/structconv"
)

func TestAddCountries(t *testing.T) {
	db, err := New(path)
	if err != nil {
		t.Fatal("can't create db instance:", err)
	}

	addr := "https://cdn.simplelocalize.io/public/v1/countries"
	resp, err := http.Get(addr)
	if err != nil {
		t.Fatalf("can't make get request to %s: %v", addr, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("bad status code for http: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("can't read response body: %v", err)
	}

	countries := &[]*storage.Country{}
	err = json.Unmarshal(body, countries)
	if err != nil {
		t.Fatalf("can't unmarshal response: %v", err)
	}

	// time.Sleep(time.Millisecond * 200)
	for _, country := range *countries {

		// Save country
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		id, err := db.SaveCountry(ctx, country)
		cancel()
		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("added Country with id %d\n", id)

	}
}

func TestAddGetRemoveCountry(t *testing.T) {
	db, err := New(path)
	if err != nil {
		t.Fatal("can't create db instance:", err)
	}

	time.Sleep(time.Millisecond * 200)

	testCoutry := &storage.Country{
		CountryName: "test",
		CountryCode: "TS",
	}

	// Save country
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	id, err := db.SaveCountry(ctx, testCoutry)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("added Country with id %d\n", id)
	time.Sleep(time.Millisecond * 100)

	// Get country
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	country, err := db.GetCountryByID(ctx, id)
	cancel()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("country: %+v", country)

	// Update country
	country.CountryName = "test-new"
	country.CountryCode = "TS-new"

	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	id, err = db.SaveCountry(ctx, country)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("added Country with id %d\n", id)
	time.Sleep(time.Millisecond * 100)

	// Test Get country
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	countryUpd, err := db.GetCountryByID(ctx, id)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	if err := structconv.CompareStructs(country, countryUpd, []string{"CountryName", "CountryCode"}); err != nil {
		t.Fatal(err)
	}

	// Delete country
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	err = db.RemoveCountryByID(ctx, id)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	// Test Get country
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	country, err = db.GetCountryByID(ctx, id)
	cancel()
	if err == storage.ErrNoSuchCountry {
		fmt.Println("Country deleted successfully")
	} else if err != nil {
		t.Fatal(err)
	} else {
		t.Fatal("Country not deleted")
	}

}
