package x_ui

import (
	"context"
	"fmt"
	"net/http"
	"vpn-tg-bot/pkg/clients/x-ui/model"
)

const (
	AllSettingsPath = SettingPrefix + "/all"
)

func (c *XUIClient) GetAllSettings(ctx context.Context) (settings *model.AllSetting, err error) {
	defer func() {
		if err == nil {
			c.Settings = settings
		}
	}()

	resp, err := c.post(ctx, AllSettingsPath, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	return assignResponseTo[model.AllSetting](resp)
}
