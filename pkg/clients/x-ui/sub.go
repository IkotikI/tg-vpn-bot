package x_ui

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	SubLinkPath = "/{{.SubPath}}/{{.SubID}}"
)

func (c *XUIClient) GetSubBySubID(ctx context.Context, subID string) (link string, err error) {

	if c.Settings == nil {
		c.GetAllSettings(ctx)
	}

	if !c.Settings.SubEnable {
		return "", errors.New("sub is disabled on remote server")
	}

	subPath := strings.Trim(c.Settings.SubPath, "/")
	subPort := c.Settings.SubPort
	subEncrypt := c.Settings.SubEncrypt

	args := struct {
		SubPath string
		SubID   string
	}{
		SubPath: subPath,
		SubID:   subID,
	}

	path, err := c.PreparePath(SubLinkPath, args)
	if err != nil {
		return "", err
	}

	u := url.URL{
		Scheme: c.Protocol,
		Host:   c.Host + ":" + strconv.Itoa(subPort),
		Path:   path,
	}

	resp, err := c.get(ctx, u.String(), nil)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("bad status code: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if subEncrypt {
		bytes, err := base64.StdEncoding.DecodeString(string(body))
		if err != nil {
			log.Printf("[ERR] clients/x-ui: GetSubBySubID: can't decode base64 string: the string is \"%s\"", string(body))
			return "", err
		}

		link = string(bytes)
	} else {
		link = string(body)
	}

	return link, nil
}
