package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/structconv"
	"vpn-tg-bot/web/admin_panel"
	"vpn-tg-bot/web/admin_panel/controller"
)

func TestUpdateUser(t *testing.T) {
	s := getTestSettings(t)
	time.Sleep(100 * time.Microsecond)

	user := &storage.User{
		TelegramID:   101,
		TelegramName: "test",
	}

	// Add test User
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	id, err := s.Storage.SaveUser(ctx, user)
	cancel()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("User added with id: %d", id)
	// Defer Delete test User
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		err = s.Storage.RemoveUserByID(ctx, id)
		cancel()
		if err != nil {
			t.Log(err)
		}
		t.Log("User removed")
	}()

	time.Sleep(100 * time.Microsecond)

	user.TelegramName = "test_new1"

	updateViaJSON(t, s, id, user)
	checkUser(t, s, id, user)

	time.Sleep(100 * time.Microsecond)

	user.TelegramName = "test_new2"
	updateViaURLEncoding(t, s, id, user)
	checkUser(t, s, id, user)

	// Check, if user have been changed.

}

func updateViaURLEncoding(t *testing.T, s *admin_panel.Settings, id storage.UserID, user *storage.User) {
	formData := url.Values{}
	formData.Set("telegram_id", fmt.Sprintf("%d", user.TelegramID))
	formData.Set("telegram_name", user.TelegramName)

	encodedFormData := formData.Encode()
	buffer := bytes.NewBufferString(encodedFormData)

	host := fmt.Sprintf("%s://%s/user/%d", s.Scheme, s.Addr, id)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	request, err := http.NewRequestWithContext(ctx, http.MethodPatch, host, buffer)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpResp, err := http.DefaultClient.Do(request)
	cancel()
	if err != nil {
		t.Fatal(err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		decoder := json.NewDecoder(httpResp.Body)
		resp := &controller.Response{}
		err := decoder.Decode(resp)
		if err != nil {
			t.Fatalf("bad status code: %d; can'te decode response", httpResp.StatusCode)
		}
		t.Fatalf("bad status code: %d, message: %s", httpResp.StatusCode, resp.Msg)
	}
	t.Log("URLEncoding request succeed")

}

func updateViaJSON(t *testing.T, s *admin_panel.Settings, id storage.UserID, user *storage.User) {

	payload, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(payload))
	buffer := bytes.NewBuffer(payload)

	host := fmt.Sprintf("%s://%s/user/%d", s.Scheme, s.Addr, id)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, host, buffer)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	httpResp, err := http.DefaultClient.Do(request)
	cancel()
	if err != nil {
		t.Fatal(err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		decoder := json.NewDecoder(httpResp.Body)
		resp := &controller.Response{}
		err := decoder.Decode(resp)
		if err != nil {
			t.Fatalf("bad status code: %d; can'te decode response", httpResp.StatusCode)
		}
		t.Fatalf("bad status code: %d, message: %s", httpResp.StatusCode, resp.Msg)
	}
	t.Log("JSON request succeed")
}

func checkUser(t *testing.T, s *admin_panel.Settings, id storage.UserID, user *storage.User) {
	// Check, if user have been changed.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	gotUser, err := s.Storage.GetUserByID(ctx, id)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	if err := structconv.CompareStructs(user, gotUser, []string{"TelegramName"}); err != nil {
		t.Fatal(err)
	}
	t.Logf("user have updated: \n%+v\n", gotUser)
}
