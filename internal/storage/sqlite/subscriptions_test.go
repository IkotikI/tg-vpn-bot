package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/structconv"
)

func ChangeSubsription(t *testing.T, userID storage.UserID, serverID storage.ServerID, subscription *storage.Subscription) {

}

func TestAddRemoveGetSubscription(t *testing.T) {
	db, err := New(path)
	if err != nil {
		t.Fatal("can't create db instance:", err)
	}

	userID := storage.UserID(100)
	serverID := storage.ServerID(100)

	time.Sleep(time.Millisecond * 200)

	subscription := &storage.Subscription{
		UserID:                userID,
		ServerID:              serverID,
		SubscriptionStatus:    "test",
		SubscriptionExpiredAt: time.Now().Add(time.Hour),
	}

	// Add Subscription
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	err = db.SaveSubscription(ctx, subscription)
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	fmt.Print("Subscription added\n\n")

	// Get Subscription by IDs
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	subscriptionByIDs, err := db.GetSubscriptionByIDs(ctx, userID, serverID)
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	fmt.Printf("Subscription by IDs: %+v\n\n", subscriptionByIDs)

	if err := structconv.CompareStructs(subscription, subscriptionByIDs, []string{"UserID", "ServerID", "SubscriptionStatus", "SubscriptionExpiredAt"}); err != nil {
		t.Fatalf("Subscriptions is not equal: %v", err)
	}

	// Change values
	subscription.SubscriptionStatus = "test_changed_status"
	subscription.SubscriptionExpiredAt = time.Now().Add(2 * time.Hour)

	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	err = db.SaveSubscription(ctx, subscription)
	if err != nil {
		t.Fatal(err)
	}
	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	subscriptionByIDs, err = db.GetSubscriptionByIDs(ctx, userID, serverID)
	if err != nil {
		t.Fatal(err)
	}
	cancel()

	if err = structconv.CompareStructs(subscription, subscriptionByIDs, []string{"UserID", "ServerID", "SubscriptionStatus", "SubscriptionExpiredAt"}); err != nil {
		t.Fatalf("Changed Subscriptions is not equal: %v", err)
	}
	fmt.Print("Subscription changed successfully\n\n")

	// Get Subscription by UserID
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	subscriptionByUserID, err := db.GetSubscriptionsByUserID(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	fmt.Printf("Subscriptions by UserID: %+v\n\n", subscriptionByUserID)

	// if err := structconv.CompareStructs(subscription, subscriptionByUserID, []string{"UserID", "ServerID", "SubscriptionStatus", "SubscriptionExpiredAt"}); err != nil {
	// 	t.Fatalf("Subscriptions is not equal: %v", err)
	// }

	// Get Subscription by ServerID
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	subscriptionByServerID, err := db.GetSubscriptionsServerID(ctx, serverID)
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	fmt.Printf("Subscriptions by ServerID: %+v\n\n", subscriptionByServerID)

	// if err := structconv.CompareStructs(subscription, subscriptionByServerID, []string{"UserID", "ServerID", "SubscriptionStatus", "SubscriptionExpiredAt"}); err != nil {
	// 	t.Fatalf("Subscriptions is not equal: %v", err)
	// }

	// Remove Subscription
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	err = db.RemoveSubscriptionByIDs(ctx, userID, serverID)
	if err != nil {
		t.Fatal(err)
	}
	cancel()

	// Get Subscription by UserID after removing
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	_, err = db.GetSubscriptionByIDs(ctx, userID, serverID)
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}
	cancel()
	fmt.Println("Subscription removed successfully")

}

func TestAddSubscriptions(t *testing.T) {
	db, err := New(path)
	if err != nil {
		t.Fatal("can't create db instance:", err)
	}

	time.Sleep(time.Millisecond * 200)

	n := 8
	for i := 1; i <= n; i++ {

		subscription := &storage.Subscription{
			UserID:                storage.UserID(i),
			ServerID:              storage.ServerID(i),
			SubscriptionStatus:    "test",
			SubscriptionExpiredAt: time.Now().Add(time.Hour),
		}

		// Add Subscription
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		err = db.SaveSubscription(ctx, subscription)
		if err != nil {
			t.Fatal(err)
		}
		cancel()
		fmt.Printf("Added subscription with id %d\n", subscription.UserID)

	}

	db.Close()
}

func TestRemoveSubscriptions(t *testing.T) {
	db, err := New(path)
	if err != nil {
		t.Fatal("can't create db instance:", err)
	}

	time.Sleep(time.Millisecond * 200)

	n := 8
	for i := 1; i <= n; i++ {

		userID := storage.UserID(i)
		serverID := storage.ServerID(i)

		// Add Subscription
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		err = db.RemoveSubscriptionByIDs(ctx, userID, serverID)
		if err != nil {
			t.Fatal(err)
		}
		cancel()
		fmt.Printf("Removed subscription with user id %d and server id %d\n", userID, serverID)

	}

	db.Close()
}
