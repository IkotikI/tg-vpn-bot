package storage

import (
	"net/url"
	"strconv"
	"time"
)

func (u *User) ParseURLValues(url url.Values, time_layout string) {
	id, err := strconv.ParseInt(url.Get("id"), 10, 64)
	if err != nil {
		id = 0
	}
	u.ID = UserID(id)

	telegram_id, err := strconv.ParseInt(url.Get("telegram_id"), 10, 64)
	if err != nil {
		telegram_id = 0
	}
	u.TelegramID = TelegramID(telegram_id)

	u.TelegramName = url.Get("telegram_name")

	created_at, err := time.Parse(time_layout, url.Get("created_at"))
	if err != nil {
		created_at = time.Time{}
	}
	u.CreatedAt = created_at

	updated_at, err := time.Parse(time_layout, url.Get("updated_at"))
	if err != nil {
		updated_at = time.Time{}
	}
	u.UpdatedAt = updated_at
}

func (u *User) ParseDefaultsFrom(from *User) {
	if u.ID == 0 {
		u.ID = from.ID
	}
	if u.TelegramID == 0 {
		u.TelegramID = from.TelegramID
	}
	if u.TelegramName == "" {
		u.TelegramName = from.TelegramName
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = from.CreatedAt
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = from.UpdatedAt
	}

}

func (s *VPNServer) ParseURLValues(url url.Values, time_layout string) {
	id, err := strconv.ParseInt(url.Get("id"), 10, 64)
	if err != nil {
		id = 0
	}
	s.ID = ServerID(id)

	country_id, err := strconv.ParseInt(url.Get("country_id"), 10, 64)
	if err != nil {
		country_id = 0
	}
	s.CountryID = CountryID(country_id)

	s.Name = url.Get("name")

	s.Protocol = url.Get("protocol")

	s.Host = url.Get("host")

	port, err := strconv.Atoi(url.Get("port"))
	if err != nil {
		port = 0
	}
	s.Port = port

	s.Username = url.Get("username")

	s.Password = url.Get("password")

	created_at, err := time.Parse(time_layout, url.Get("created_at"))
	if err != nil {
		created_at = time.Time{}
	}
	s.CreatedAt = created_at

	updated_at, err := time.Parse(time_layout, url.Get("updated_at"))
	if err != nil {
		updated_at = time.Time{}
	}
	s.UpdatedAt = updated_at
}

func (s *VPNServer) ParseDefaultsFrom(from *VPNServer) {
	if s.ID == 0 {
		s.ID = from.ID
	}
	if s.CountryID == 0 {
		s.CountryID = from.CountryID
	}
	if s.Protocol == "" {
		s.Protocol = from.Protocol
	}
	if s.Host == "" {
		s.Host = from.Host
	}
	if s.Port == 0 {
		s.Port = from.Port
	}
	if s.Username == "" {
		s.Username = from.Username
	}
	if s.Password == "" {
		s.Password = from.Password
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = from.CreatedAt
	}
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = from.UpdatedAt
	}
}

func (s *Subscription) ParseDefaultsFrom(from *Subscription) {
	if s.UserID == 0 {
		s.UserID = from.UserID
	}
	if s.ServerID == 0 {
		s.ServerID = from.ServerID
	}
	if s.SubscriptionStatus == "" {
		s.SubscriptionStatus = from.SubscriptionStatus
	}
	if s.SubscriptionExpiredAt.IsZero() {
		s.SubscriptionExpiredAt = from.SubscriptionExpiredAt
	}
}

func (s *Subscription) ParseURLValues(url url.Values, time_layout string) {
	user_id, err := strconv.ParseInt(url.Get("user_id"), 10, 64)
	if err != nil {
		user_id = 0
	}
	s.UserID = UserID(user_id)

	server_id, err := strconv.ParseInt(url.Get("server_id"), 10, 64)
	if err != nil {
		server_id = 0
	}
	s.ServerID = ServerID(server_id)

	s.SubscriptionStatus = url.Get("subscription_status")

	s.SubscriptionExpiredAt, err = time.Parse(time_layout, url.Get("subscription_expired_at"))
	if err != nil {
		s.SubscriptionExpiredAt = time.Time{}
	}
}

func (s *Subscription) IsCorrectStatus() bool {
	switch s.SubscriptionStatus {
	case SubscriptionStatusActive, SubscriptionStatusExpired, SubscriptionStatusDisabled:
		return true
	}
	return false
}
