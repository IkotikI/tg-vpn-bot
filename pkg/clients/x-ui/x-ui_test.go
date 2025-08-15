package x_ui

import (
	"context"
	"net/http"
	"testing"
	"time"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/internal/storage/sqlite"
	"vpn-tg-bot/pkg/clients/x-ui/model"
	"vpn-tg-bot/pkg/e"
)

var testStoragePath = "../../../internal/storage/sqlite/test_data/db.db"

var server = &storage.VPNServer{
	ID:        1,
	CountryID: 1,
	Name:      "Local",
	Protocol:  "http",
	Host:      "localhost",
	Port:      2053,
	Username:  "admin",
	Password:  "admin",
}

func makeAuthStore() storage.ServerAuthorizations {
	db, err := sqlite.New(testStoragePath)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Millisecond * 200)
	return db
}

func TestFunctions(t *testing.T) {
	xui := New(TokenKey_3x_ui, server, makeAuthStore())
	t.Run("ParseClientID", func(t *testing.T) {
		id, err := ParseClientID(test_uuid)
		if err != nil {
			t.Fatal(err)
		}
		if test_uuid != id.String() {
			t.Fatalf("want: %s, got: %s", test_uuid, id.String())
		}
	})
	t.Run("authFromLoginCookies", func(t *testing.T) {
		cookiesExampleContent := "3x-ui=MTczNTc0MjIzMHxEWDhFQVFMX2dBQUJFQUVRQUFCMV80QUFBUVp6ZEhKcGJtY01EQUFLVEU5SFNVNWZWVk5GVWhoNExYVnBMMlJoZEdGaVlYTmxMMjF2WkdWc0xsVnpaWExfZ1FNQkFRUlZjMlZ5QWYtQ0FBRUVBUUpKWkFFRUFBRUlWWE5sY201aGJXVUJEQUFCQ0ZCaGMzTjNiM0prQVF3QUFRdE1iMmRwYmxObFkzSmxkQUVNQUFBQUZQLUNFUUVDQVFWaFpHMXBiZ0VGWVdSdGFXNEF8PwWtcxJNuxiUDtstKdZCPxm3QW-KH7gAMF-KM8qz9sI=; Path=/; Expires=Wed, 01 Jan 2025 15:37:10 GMT; Max-Age=3600; HttpOnly"
		cookiesExample, err := http.ParseSetCookie(cookiesExampleContent)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("cookie: %+v\n", cookiesExample)

		auth, err := xui.parseAuthCookie(cookiesExample)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("auth: %+v\n", auth)
	})

	t.Run("assignTo", func(t *testing.T) {
		xui = New(TokenKey_3x_ui, server, makeAuthStore())
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)

		args := struct{ InboundID int }{InboundID: test_inboundID}
		path, err := xui.PreparePath(GetInboundPath, args)
		if err != nil {
			t.Fatal(e.Wrap("can't prepare path", err))
		}

		// Request 1
		resp, err := xui.get(ctx, path, nil)
		if err != nil {
			t.Fatal(e.Wrap("can't send get inbound request", err))
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("bad status code: %s", resp.Status)
		}
		inbound_assign, err := assignResponseTo[model.Inbound](resp)
		if err != nil {
			t.Fatalf("can't assign inbound via generic: %s", err)
		}

		resp.Body.Close()
		cancel()
		time.Sleep(400 * time.Millisecond)

		// Request 2
		ctx, cancel = context.WithTimeout(context.Background(), 400*time.Millisecond)
		resp, err = xui.get(ctx, path, nil)
		if err != nil {
			t.Fatal(e.Wrap("can't send get inbound request", err))
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("bad status code: %s", resp.Status)
		}

		inbound_classic, err := xui.assignInbound(resp)
		if err != nil {
			t.Fatalf("can't assign inbound via classic: %s", err)
		}

		resp.Body.Close()
		cancel()

		// Compare
		if *inbound_assign != *inbound_classic {
			t.Fatalf("inbounds are not equal: \n classic: %+v \n generic: %+v", inbound_classic, inbound_assign)
		}

		t.Log(inbound_assign)
	})
}

func TestRequests(t *testing.T) {
	xui := New(TokenKey_3x_ui, server, makeAuthStore())

	t.Run("auth", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		err := xui.Auth(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})

	var inbound *model.Inbound = &DefaultInbound
	inbound.Port = inbound.Port + 1111
	t.Run("post", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		addedInbound, err := xui.AddInbound(ctx, inbound)
		if err != nil {
			t.Error(err)
		}

		t.Logf("post: Added Inbound: \n %+v \n", addedInbound)

		ctx, cancel = context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		updInbound, err := xui.UpdateInbound(ctx, addedInbound)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("post: Update Inbound: \n %+v \n", updInbound)

		ctx, cancel = context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		err = xui.DeleteInbound(ctx, addedInbound.Id)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("post: Deleted Inbound: \n %+v \n", updInbound)
	})

	t.Run("get", func(t *testing.T) {
		var err error
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		inbound, err = xui.GetInbound(ctx, 1)
		if err == nil {
			t.Fatal(err)
		}

		t.Logf("get: Get Inbound: \n %+v \n", inbound)
	})

}
