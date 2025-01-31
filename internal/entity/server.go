package entity

import "time"

type ServerID int64

type VPNServer struct {
	ID        ServerID
	Name      string
	Protocol  string
	IPaddress string
	Port      int
	Login     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time

	Country Country
}
