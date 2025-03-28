package vpnserver

import (
	"context"
	"vpn-tg-bot/internal/storage"
)

type VPNServerManager struct {
	storage   storage.Storage
	authStore storage.ServerAuthorizations
}

func New(storage storage.Storage, authStore storage.ServerAuthorizations) *VPNServerManager {
	return &VPNServerManager{storage: storage, authStore: authStore}
}

func (m *VPNServerManager) GetDemoServer(ctx context.Context) (*storage.VPNServer, error) {
	serverID := storage.ServerID(1)
	return m.storage.GetServerByID(ctx, serverID)
}
