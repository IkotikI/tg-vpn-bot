package x_ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"vpn-tg-bot/pkg/clients/x-ui/model"
	"vpn-tg-bot/pkg/e"
)

// Exclude Inbound from the response. Reduce some boilerplate.
func (c *XUIClient) assignInbound(resp *http.Response) (inbound *model.Inbound, err error) {
	defer func() { e.WrapIfErr("can't assign inbound", err) }()
	inbound = &model.Inbound{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("can't read response body", err)
	}

	respStruct := &Response{Obj: &model.Inbound{}}
	err = json.Unmarshal(body, respStruct)
	if err != nil {
		return nil, e.Wrap("can't unmarshal inbound", err)
	}

	if respStruct.Success == false {
		return nil, fmt.Errorf("server responded with error: %s", respStruct.Msg)
	}

	inbound, ok := respStruct.Obj.(*model.Inbound)
	if !ok {
		fmt.Printf("respStruct %v", respStruct)
		return nil, errors.New("can't cast Obj to Inbound")
	}

	return inbound, nil
}
