package x_ui

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"vpn-tg-bot/pkg/clients/x-ui/model"
	"vpn-tg-bot/pkg/e"
)

const (
	// Inbound
	AddInboundPath    = InboundsPrefix + "/add"
	GetInboundPath    = InboundsPrefix + "/get/{{.InboundID}}"
	UpdateInboundPath = InboundsPrefix + "/update/{{.InboundID}}"
	DeleteInboundPath = InboundsPrefix + "/del/{{.InboundID}}"
	// Inbounds
	GetInboundsPath = InboundsPrefix + "list"
)

// var (
// 	// Inbound
// 	ErrInboundAlreadyExists = errors.New("inbound already exists")
// 	ErrInboundNotFound      = errors.New("inbound not found")
// )

var DefaultInbound model.Inbound = model.Inbound{
	Up:             0,
	Down:           0,
	Total:          0,
	Remark:         "",
	Enable:         true,
	ExpiryTime:     0,
	Listen:         "",
	Port:           55421,
	Protocol:       model.VLESS,
	Settings:       "",
	StreamSettings: `"{\"network\": \"tcp\",\"security\": \"reality\",\"externalProxy\": [],\"realitySettings\": {\"show\": false,\"xver\": 0,\"dest\": \"yahoo.com:443\",\"serverNames\": [\"yahoo.com\",\"www.yahoo.com\"],\"privateKey\": \"wIc7zBUiTXBGxM7S7wl0nCZ663OAvzTDNqS7-bsxV3A\",\"minClient\": \"\",\"maxClient\": \"\",\"maxTimediff\": 0,\"shortIds\": [\"47595474\",\"7a5e30\",\"810c1efd750030e8\",\"99\",\"9c19c134b8\",\"35fd\",\"2409c639a707b4\",\"c98fc6b39f45\"],\"settings\": {\"publicKey\": \"2UqLjQFhlvLcY7VzaKRotIDQFOgAJe1dYD1njigp9wk\",\"fingerprint\": \"random\",\"serverName\": \"\",\"spiderX\": \"/\"}},\"tcpSettings\": {\"acceptProxyProtocol\": false,\"header\": {\"type\": \"none\"}}}"`,
	Sniffing:       `"{\"enabled\": true,\"destOverride\": [\"http\",\"tls\",\"quic\",\"fakedns\"],\"metadataOnly\": false,\"routeOnly\": false}"`,
	Allocate:       `"{\"strategy\": \"always\",\"refresh\": 5,\"concurrency\": 3}"`,
}

/* ---- Inbound ---- */
func (c *XUIClient) AddInbound(ctx context.Context, inbound *model.Inbound) (addedInbound *model.Inbound, err error) {
	defer func() { err = e.WrapIfErr("can't add inbound", err) }()

	data, err := json.Marshal(inbound)
	if err != nil {
		return nil, e.Wrap("can't marshal inbound", err)
	}

	buffer := bytes.NewBuffer(data)
	resp, err := c.post(ctx, AddInboundPath, buffer)
	if err != nil {
		return nil, e.Wrap("can't send add inbound request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %s", resp.Status)
	}

	return assignResponseTo[model.Inbound](resp)
}

func (c *XUIClient) UpdateInbound(ctx context.Context, inbound *model.Inbound) (updatedInbound *model.Inbound, err error) {
	defer func() { err = e.WrapIfErr("can't update inbound", err) }()
	if inbound == nil {
		return nil, errors.New("inbound is nil")
	}
	if inbound.Id == 0 {
		return nil, errors.New("inbound id is zero")
	}

	data, err := json.Marshal(inbound)
	if err != nil {
		return nil, e.Wrap("can't marshal inbound", err)
	}

	args := struct{ InboundID int }{InboundID: inbound.Id}
	path, err := c.PreparePath(UpdateInboundPath, args)
	if err != nil {
		return nil, e.Wrap("can't prepare path", err)
	}

	buffer := bytes.NewBuffer(data)
	resp, err := c.post(ctx, path, buffer)
	if err != nil {
		return nil, e.Wrap("can't send add inbound request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %s", resp.Status)
	}

	return assignResponseTo[model.Inbound](resp)
}

func (c *XUIClient) DeleteInbound(ctx context.Context, inboundID int) (err error) {
	defer func() { err = e.WrapIfErr("can't delete inbound", err) }()

	args := struct{ InboundID int }{InboundID: inboundID}
	path, err := c.PreparePath(DeleteInboundPath, args)
	if err != nil {
		return e.Wrap("can't prepare path", err)
	}

	resp, err := c.post(ctx, path, nil)
	if err != nil {
		return e.Wrap("can't send delete inbound request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return e.Wrap("can't read response body", err)
	}

	respStruct := &Response{Obj: &model.Inbound{}}
	err = json.Unmarshal(body, respStruct)
	if err != nil {
		return e.Wrap("can't unmarshal inbound", err)
	}

	if !respStruct.Success {
		return fmt.Errorf("server responded with error: %s", respStruct.Msg)
	}

	return nil
}

func (c *XUIClient) GetInbound(ctx context.Context, inboundID int) (inbound *model.Inbound, err error) {
	defer func() { err = e.WrapIfErr("can't get inbound", err) }()

	args := struct{ InboundID int }{InboundID: inboundID}
	path, err := c.PreparePath(GetInboundPath, args)
	if err != nil {
		return nil, e.Wrap("can't prepare path", err)
	}

	resp, err := c.get(ctx, path, nil)
	if err != nil {
		return nil, e.Wrap("can't send get inbound request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %s", resp.Status)
	}

	return assignResponseTo[model.Inbound](resp)
}
