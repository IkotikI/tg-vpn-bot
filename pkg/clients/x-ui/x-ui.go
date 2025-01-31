package x_ui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/clients/x-ui/model"
	"vpn-tg-bot/pkg/e"
)

type XUIClient struct {
	serverID storage.ServerID
	Protocol string
	Host     string
	Port     int
	Username string
	password string
	token    string
	TokenKey string

	MaxRetry   int
	httpClient *http.Client
	authStore  storage.ServerAuthorizations
}

type Response struct {
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
	Obj     interface{} `json:"obj"`
}

const (
	// Login
	LoginPath = "/login"
	// Client
	InboundsPrefix = "/panel/api/inbounds"
	//Onlines
	OnlinesPath = InboundsPrefix + "/onlines"
)

const (
	TokenKey_3x_ui = "3x-ui"
	TokenKey_x_ui  = "x-ui"
)

// New returns new XUIClient. TokenKey define, which key will be used to parse auth token from cookies.
func New(tokenKey string, server storage.VPNServer, authStore storage.ServerAuthorizations) *XUIClient {
	return &XUIClient{
		serverID: server.ID,
		Protocol: server.Protocol,
		Host:     server.Host,
		Port:     server.Port,
		Username: server.Username,
		password: server.Password,
		TokenKey: tokenKey,
		MaxRetry: 3,

		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		authStore: authStore,
	}
}

func (c *XUIClient) Onlines(ctx context.Context) (emails *[]string, err error) {
	defer func() { e.WrapIfErr("can't get onlines", err) }()

	resp, err := c.post(ctx, OnlinesPath, nil)
	if err != nil {
		return nil, e.Wrap("can't get onlines", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %s", resp.Status)
	}

	return assignResponseTo[[]string](resp)
}

func (c *XUIClient) PreparePath(path string, args interface{}) (activePath string, err error) {
	temp, err := template.New(path).Parse(path)
	if err != nil {
		return "", err
	}

	b := &bytes.Buffer{}
	err = temp.Execute(b, args)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func (c *XUIClient) APILinkURL(path string) (u *url.URL) {
	if strings.Count(path, "://") > 0 {
		u, _ = url.Parse(path)
	} else {
		u = &url.URL{
			Scheme: c.Protocol,
			Host:   c.Host + ":" + strconv.Itoa(c.Port),
			Path:   path,
		}
	}
	return u
}

func (c *XUIClient) APILink(path string) string {
	return c.APILinkURL(path).String()
}

func (c *XUIClient) get(ctx context.Context, path string, body io.Reader) (httpResp *http.Response, err error) {
	defer func() { e.WrapIfErr("can't make request", err) }()

	for retry := 0; retry < c.MaxRetry; retry++ {

		if c.token == "" {
			err = c.Auth(ctx)
			if err != nil {
				log.Printf("can't login 3x-ui, retry %d, error: %s", retry, err)
			}
		}

		authCookie := http.Cookie{Name: c.TokenKey, Value: c.token}

		fmt.Println("c.APILink(path)", c.APILink(path))
		request, err := http.NewRequestWithContext(ctx, "GET", c.APILink(path), body)
		if err != nil {
			// Error making request depends of bad internal logic. Retries useless.
			return nil, e.Wrap("can't create get request", err)
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Cookie", authCookie.String())

		httpResp, err = c.httpClient.Do(request)
		if err != nil {
			log.Printf("can't execute get request, retry %d, error: %s", retry, err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Debug
		fmt.Println("post response content type", httpResp.Header.Get("Content-Type"))
		fmt.Printf("full api path: %s\n", c.APILink(path))
		// fmt.Printf("cookies: %s, from %+v\n", authCookie.String(), authCookie)
		// Debug End

		if httpResp.StatusCode != http.StatusOK || httpResp.Header.Get("Content-Type") != "application/json" {
			err = c.login(ctx)
			if err != nil {
				if retry != 0 {
					log.Printf("can't login 3x-ui, retry %d, error: %s", retry, err)
				}
				time.Sleep(100 * time.Millisecond)
				continue
			}
		}

		return httpResp, nil

	}

	return nil, err
}

func (c *XUIClient) post(ctx context.Context, path string, body io.Reader) (httpResp *http.Response, err error) {
	defer func() { e.WrapIfErr("can't make request", err) }()

	for retry := 0; retry < c.MaxRetry; retry++ {

		if c.token == "" {
			err = c.Auth(ctx)
			if err != nil {
				log.Printf("can't login 3x-ui, retry %d, error: %s", retry, err)
			}
		}

		authCookie := http.Cookie{Name: c.TokenKey, Value: c.token}

		request, err := http.NewRequestWithContext(ctx, "POST", c.APILink(path), body)
		if err != nil {
			// Error making request depends of bad internal logic. Retries useless.
			return nil, e.Wrap("can't create post request", err)
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Cookie", authCookie.String())

		httpResp, err = c.httpClient.Do(request)
		if err != nil {
			log.Printf("can't execute post request, retry %d, error: %s", retry, err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Debug
		fmt.Println("post response content type", httpResp.Header.Get("Content-Type"))
		fmt.Printf("full api path: %s\n", c.APILink(path))
		// fmt.Printf("cookies: %s, from %+v\n", authCookie.String(), authCookie)
		// Debug End

		if httpResp.StatusCode != http.StatusOK || httpResp.Header.Get("Content-Type") != "application/json" {
			err = c.login(ctx)
			if err != nil {
				if retry != 0 {
					log.Printf("can't login 3x-ui, retry %d, error: %s", retry, err)
				}
				time.Sleep(100 * time.Millisecond)
				continue
			}
		}

		return httpResp, nil

	}

	return nil, err
}

func (c *XUIClient) login(ctx context.Context) (err error) {
	defer func() { e.WrapIfErr("can't login 3x-ui", err) }()
	fmt.Println("try to login")

	payload := model.LoginRequest{
		Username: c.Username,
		Password: c.password,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return e.Wrap("can't marshal login request", err)
	}

	buffer := bytes.NewBuffer(data)
	resp, err := http.Post(c.APILink(LoginPath), "application/json", buffer)
	if err != nil {
		return e.Wrap("can't send login request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return e.Wrap("can't read response body", err)
	}

	respObj := &Response{}
	err = json.Unmarshal(body, respObj)
	if err != nil {
		return e.Wrap("can't unmarshal login response", err)
	}
	if !respObj.Success {
		return fmt.Errorf("login failed: %s", respObj.Msg)
	}

	auth, err := c.authFromLoginCookies(resp.Cookies())
	if err != nil {
		return e.Wrap("can't parse auth cookies", err)
	}

	// fmt.Printf("auth data %+v\n\n", auth)
	_, err = c.authStore.SaveServerAuth(ctx, auth)
	if err != nil {
		return e.Wrap("can't save authorization data", err)
	}

	c.token = auth.Token

	return nil
}

func (c *XUIClient) Auth(ctx context.Context) error {
	auth, err := c.authStore.GetServerAuthByServerID(ctx, c.serverID)
	if err == storage.ErrNoSuchServerAuth {
		fmt.Println("auth not found in store")
		if err = c.login(ctx); err != nil {
			return err
		}
		if _, err = c.Onlines(ctx); err != nil {
			return err
		}
	} else if err != nil {
		return e.Wrap("can't get authorization data", err)
	} else {
		if auth.ExpiredAt.Before(time.Now()) {
			fmt.Println("auth token has expired")
			if err = c.login(ctx); err != nil {
				return err
			}
		} else {
			c.token = auth.Token
		}
	}
	// fmt.Printf("auth %+v\n", auth)
	// fmt.Printf("token %s\n", c.token)

	return nil
}

func (c *XUIClient) parseAuthCookie(cookie *http.Cookie) (auth *storage.VPNServerAuthorization, err error) {
	auth = &storage.VPNServerAuthorization{}

	auth.Token = cookie.Value
	auth.ExpiredAt = cookie.Expires
	auth.Meta = cookie.String()

	return auth, nil
}

func (c *XUIClient) authFromLoginCookies(cookies []*http.Cookie) (auth *storage.VPNServerAuthorization, err error) {

	for _, cookie := range cookies {
		switch cookie.Name {
		case c.TokenKey:
			auth, err := c.parseAuthCookie(cookie)
			if err != nil {
				return nil, e.Wrap("can't parse auth cookie", err)
			}
			auth.ServerID = c.serverID
			return auth, nil
		}
	}

	return nil, fmt.Errorf("token by cookie field '%s' don't found", c.TokenKey)
}

func (c *XUIClient) ServerID() storage.ServerID {
	return c.serverID
}

// Unmarshal http.Response.Body and assign data from .Obj field to given type T.
// See: Response{} type; model/model.go.
func assignResponseTo[T any](resp *http.Response) (t *T, err error) {
	defer func() { e.WrapIfErr(fmt.Sprintf("can't assign %T", t), err) }()
	t = new(T)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("can't read response body", err)
	}

	respStruct := &Response{Obj: new(T)}
	err = json.Unmarshal(body, respStruct)
	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("can't unmarshal %T", t), err)
	}

	if respStruct.Success == false {
		return nil, fmt.Errorf("server responded with error: \"%s\"", respStruct.Msg)
	}

	t, ok := respStruct.Obj.(*T)
	fmt.Printf("respStruct.Obj.(%T) %+v", t, t)
	if !ok {
		fmt.Printf("respStruct %+v", respStruct)
		return nil, fmt.Errorf("can't cast Obj to %T", t)
	}
	if t == nil {
		return nil, fmt.Errorf("nil pointer, instead of %T", t)
	}

	return t, nil
}
