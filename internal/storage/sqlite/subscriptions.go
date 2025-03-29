package sqlite

import (
	"context"
	"database/sql"
	"vpn-tg-bot/internal/storage"
)

/* ---- Interface implementation ---- */

func (s *SQLStorage) GetSubscriptions(ctx context.Context, args *storage.QueryArgs) (subs *[]storage.Subscription, err error) {
	q := `
		SELECT * FROM subscriptions
	`
	queryEnd, queryArgs := s.buildParts([]string{"where", "order_by", "limit"}, args)
	q += queryEnd

	subs = &[]storage.Subscription{}
	err = s.db.SelectContext(ctx, subs, q, queryArgs...)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

/* ---- Writer Interface impoelementation ---- */

func (s *SQLStorage) SaveSubscription(ctx context.Context, sub *storage.Subscription) (err error) {
	q := `
		INSERT INTO subscriptions (user_id, server_id, subscription_status, subscription_expired_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(user_id, server_id)
		DO UPDATE SET 
		subscription_status = excluded.subscription_status,
		subscription_expired_at = excluded.subscription_expired_at
	`

	_, err = s.db.ExecContext(ctx, q, sub.UserID, sub.ServerID, sub.SubscriptionStatus, sub.SubscriptionExpiredAt)

	// oldSub, err := s.GetSubscriptionByIDs(ctx, sub.UserID, sub.ServerID)

	// if err == sql.ErrNoRows {
	// 	q := `INSERT INTO subscriptions (user_id, server_id, subscription_status, subscription_expired_at) VALUES (?, ?, ?, ?)`
	// 	_, err = s.db.ExecContext(ctx, q, sub.UserID, sub.ServerID, sub.SubscriptionStatus, sub.SubscriptionExpiredAt)
	// 	if err != nil {
	// 		return e.Wrap("can't execute query", err)
	// 	}
	// } else if err != nil {
	// 	return e.Wrap("can't scan row", err)
	// } else {
	// 	q := `UPDATE subscriptions SET subscription_status = ?, subscription_expired_at = ? WHERE user_id = ? AND server_id = ?`
	// 	_, err = s.db.ExecContext(ctx, q, sub.SubscriptionStatus, sub.SubscriptionExpiredAt, sub.UserID, sub.ServerID)
	// 	if err != nil {
	// 		return e.Wrap("can't execute query", err)
	// 	}
	// }

	return err
}

func (s *SQLStorage) RemoveSubscriptionByIDs(ctx context.Context, userID storage.UserID, serverID storage.ServerID) (err error) {
	q := `DELETE FROM subscriptions WHERE user_id = ? AND server_id = ?`

	_, err = s.db.ExecContext(ctx, q, userID, serverID)
	return err
}

/* ---- Reader Interface implementation ---- */

func (s *SQLStorage) GetSubscriptionByIDs(ctx context.Context, userID storage.UserID, serverID storage.ServerID) (subscription *storage.Subscription, err error) {
	q := `SELECT * FROM subscriptions WHERE user_id = ? AND server_id = ? LIMIT 1`

	subscription = &storage.Subscription{}
	err = s.db.GetContext(ctx, subscription, q, userID, serverID)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSuchSubscription
	}
	return subscription, err
}

func (s *SQLStorage) GetSubscriptionsByUserID(ctx context.Context, userID storage.UserID) (subscriptions *[]storage.Subscription, err error) {
	q := `SELECT * FROM subscriptions WHERE user_id = ?`

	subscriptions = &[]storage.Subscription{}
	err = s.db.SelectContext(ctx, subscriptions, q, userID)
	return subscriptions, err
}

func (s *SQLStorage) GetSubscriptionsServerID(ctx context.Context, serverID storage.ServerID) (subscriptions *[]storage.Subscription, err error) {
	q := `SELECT * FROM subscriptions WHERE server_id = ?`

	subscriptions = &[]storage.Subscription{}
	err = s.db.SelectContext(ctx, subscriptions, q, serverID)
	return subscriptions, err
}
