package sqlite

import (
	"context"
	"vpn-tg-bot/internal/storage"
)

/* ---- Interface implementation ---- */

func (s *SQLStorage) UpdateSubscription(ctx context.Context, sub *storage.Subscription) error {
	q := "INSERT INTO subscriptions (user_id, server_id, status, expired_at) VALUES (?, ?, ?, ?)"
	_, err := s.db.ExecContext(ctx, q, sub.UserID, sub.ServerID, sub.SubscriptionStatus, sub.SubscriptionExpiredAt)
	return err
}

func (s *SQLStorage) RemoveSubscriptionByID(ctx context.Context, userID storage.UserID, serverID storage.ServerID) error {
	q := `DELETE FROM subscriptions WHERE user_id = ? AND server_id = ?`
	_, err := s.db.ExecContext(ctx, q, userID, serverID)
	return err
}
