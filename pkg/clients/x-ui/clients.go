package x_ui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"vpn-tg-bot/pkg/clients/x-ui/model"
	"vpn-tg-bot/pkg/e"

	"github.com/google/uuid"
)

const (
	// Clients
	AddClientPath        = InboundsPrefix + "/addClient"
	GetClientByIDPath    = InboundsPrefix + "/getClientTrafficsById/{{.ClientID}}"
	GetClientByEmailPath = InboundsPrefix + "/getClientTraffics/{{.Email}}"
	UpdateClientPath     = InboundsPrefix + "/updateClient/{{.ClientID}}"
	DeleteClientPath     = InboundsPrefix + "/{{.InboundID}}/delClient/{{.ClientID}}"
)

// var (
//
//	// Clients
//	ErrClientAlreadyExists = errors.New("client already exists")
//	ErrClientNotFound      = errors.New("client not found")
//
// )
type ClientID = uuid.UUID

func ParseClientID(s string) (cID ClientID, err error) {
	id, err := uuid.Parse(s)
	return ClientID(id), err
}

var ClientIDNil ClientID = ClientID(uuid.Nil)

type ClientRequest struct {
	InboundID int    `json:"id"`
	Settings  string `json:"settings"`
}

type Settings struct {
	Clients []*model.Client `json:"clients"`
}

/* ---- Client ---- */
func (c *XUIClient) AddClient(ctx context.Context, inboundID int, client *model.Client) (err error) {
	payload, err := c.prepareClientPayload(inboundID, client)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(payload)
	resp, err := c.post(ctx, AddClientPath, buffer)
	if err != nil {
		return e.Wrap("can't send addClient request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return e.Wrap("can't read response body", err)
	}

	respStruct := &Response{}
	err = json.Unmarshal(body, respStruct)
	if err != nil {
		return e.Wrap("can't unmarshal response", err)
	}

	return CheckResponseError(respStruct)
}

func (c *XUIClient) UpdateClient(ctx context.Context, inboundID int, client *model.Client) (err error) {
	payload, err := c.prepareClientPayload(inboundID, client)
	if err != nil {
		return err
	}

	clientUUID, err := uuid.Parse(client.ID)
	if err != nil {
		return e.Wrap("can't parse clientID", err)
	}

	args := struct {
		ClientID ClientID
	}{
		ClientID: ClientID(clientUUID),
	}
	path, err := c.PreparePath(UpdateClientPath, args)
	if err != nil {
		return e.Wrap("can't prepare path", err)
	}

	buffer := bytes.NewBuffer(payload)
	resp, err := c.post(ctx, path, buffer)
	if err != nil {
		return e.Wrap("can't send addClient request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return e.Wrap("can't read response body", err)
	}

	respStruct := &Response{}
	err = json.Unmarshal(body, respStruct)
	if err != nil {
		return e.Wrap("can't unmarshal response", err)
	}

	return CheckResponseError(respStruct)
}

func (c *XUIClient) DeleteClient(ctx context.Context, inboundID int, clientID ClientID) (err error) {
	args := struct {
		InboundID int
		ClientID  ClientID
	}{
		InboundID: inboundID,
		ClientID:  clientID,
	}
	path, err := c.PreparePath(DeleteClientPath, args)
	if err != nil {
		return e.Wrap("can't prepare path", err)
	}

	resp, err := c.post(ctx, path, nil)
	if err != nil {
		return e.Wrap("can't send delete client request", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return e.Wrap("can't read response body", err)
	}

	respStruct := &Response{}
	err = json.Unmarshal(body, respStruct)
	if err != nil {
		return e.Wrap("can't unmarshal response", err)
	}

	return CheckResponseError(respStruct)
}

func (c *XUIClient) GetClientTrafficByID(ctx context.Context, clientID ClientID) (client *[]model.ClientTraffic, err error) {
	args := struct{ ClientID ClientID }{ClientID: clientID}
	path, err := c.PreparePath(GetClientByIDPath, args)
	if err != nil {
		return nil, e.Wrap("can't prepare path", err)
	}

	resp, err := c.get(ctx, path, nil)
	if err != nil {
		return nil, e.Wrap("can't send get client by email request", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %s", resp.Status)
	}

	return assignResponseTo[[]model.ClientTraffic](resp)
}

func (c *XUIClient) GetClientClientTrafficsByEmail(ctx context.Context, email string) (clients *model.ClientTraffic, err error) {
	args := struct{ Email string }{Email: email}
	path, err := c.PreparePath(GetClientByEmailPath, args)
	if err != nil {
		return nil, e.Wrap("can't prepare path", err)
	}

	resp, err := c.get(ctx, path, nil)
	if err != nil {
		return nil, e.Wrap("can't send get client by email request", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %s", resp.Status)
	}

	return assignResponseTo[model.ClientTraffic](resp)
}

func (c *XUIClient) prepareClientPayload(inboundID int, client *model.Client) (payload []byte, err error) {

	settings := &Settings{
		Clients: []*model.Client{client},
	}
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return nil, e.Wrap("can't marshal settings", err)
	}

	clientRequest := &ClientRequest{
		InboundID: inboundID,
		Settings:  string(settingsJSON),
	}
	payload, err = json.Marshal(clientRequest)
	if err != nil {
		return nil, e.Wrap("can't marshal client request", err)
	}

	return payload, nil
}
