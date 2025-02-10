package structconv

import (
	"fmt"
	"hash/fnv"
	"reflect"
	"testing"
	"time"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/debug"
	"vpn-tg-bot/web/admin_panel/entity"
)

func TestFunctions1(t *testing.T) {

	user := &storage.User{
		ID:           1,
		TelegramName: "test",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	nullStructType := CreateSQLNullStructType(*user)
	nullStruct := reflect.New(nullStructType).Interface()

	fmt.Printf("Original struct: %+v\n", user)
	fmt.Printf("Null struct type: %+v\n", nullStructType)
	fmt.Printf("Null struct: %+v\n", nullStruct)

	val := reflect.ValueOf(nullStruct).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fmt.Printf("Field: %s, Type: %s\n", field.Name, field.Type)
	}

	basicStruct := &storage.User{}
	ConvertSQLNullStructToBasic(nullStruct, basicStruct)

	fmt.Printf("Converted back: %+v", basicStruct)

	if reflect.TypeOf(*user) != reflect.TypeOf(*basicStruct) {
		t.Fatalf("Start and finish types are not equal. Want %+v\n, got %+v\n", reflect.TypeOf(*user), reflect.TypeOf(*basicStruct))
	}

}
func TestFunctions2(t *testing.T) {

	userWithSubscription := &entity.UserWithSubscription{
		User: storage.User{
			ID:           1,
			TelegramName: "test",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		Subscription: storage.Subscription{
			UserID:                1,
			ServerID:              1,
			SubscriptionStatus:    "active",
			SubscriptionExpiredAt: time.Now(),
		},
	}

	nullStructType := CreateSQLNullStructType(*userWithSubscription)
	nullStruct := reflect.New(nullStructType).Interface()
	val := reflect.ValueOf(nullStruct).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fmt.Printf("Field: %s, Type: %s\n", field.Name, field.Type)
	}

	basicStruct := &entity.UserWithSubscription{}
	ConvertSQLNullStructToBasic(nullStruct, basicStruct)

	if reflect.TypeOf(*userWithSubscription) != reflect.TypeOf(*basicStruct) {
		t.Fatalf("Start and finish types are not equal. Want %+v\n, got %+v\n", reflect.TypeOf(*userWithSubscription), reflect.TypeOf(*basicStruct))
	}

	fmt.Printf("Converted back: %+v", basicStruct)

}

func TestParseDefaultsStrict(t *testing.T) {

	userWithSubscriptionDefault := &entity.UserWithSubscription{
		User: storage.User{
			ID:           1,
			TelegramID:   storage.TelegramID(hashStringToInt64("test")),
			TelegramName: "test",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		Subscription: storage.Subscription{
			UserID:                1,
			ServerID:              1,
			SubscriptionStatus:    "active",
			SubscriptionExpiredAt: time.Now(),
		},
	}

	userWithSubscription := &entity.UserWithSubscription{
		User: storage.User{
			TelegramID:   storage.TelegramID(hashStringToInt64("test_new")),
			TelegramName: "test_new",
			CreatedAt:    time.Now().Add(1 * time.Hour),
			UpdatedAt:    time.Now().Add(1 * time.Hour),
		},
		Subscription: storage.Subscription{
			SubscriptionStatus:    "inactive",
			SubscriptionExpiredAt: time.Now().Add(1 * time.Hour),
		},
	}

	// Deep copy
	u := *userWithSubscriptionDefault
	u.TelegramID = userWithSubscription.TelegramID
	u.TelegramName = userWithSubscription.TelegramName
	u.CreatedAt = userWithSubscription.CreatedAt
	u.UpdatedAt = userWithSubscription.UpdatedAt
	u.SubscriptionStatus = userWithSubscription.SubscriptionStatus
	u.SubscriptionExpiredAt = userWithSubscription.SubscriptionExpiredAt

	userWithSubscriptionCheck := &u

	t.Logf("start\n%s\n", debug.JSON(userWithSubscription))
	// ParseDefaults(userWithSubscription, userWithSubscriptionDefault)
	err := ParseDefaultsStrict(userWithSubscription, userWithSubscriptionDefault)
	if err != nil {
		t.Fatalf("can't parse defaults: %s", err)
	}

	t.Logf("result\n%s\n", debug.JSON(userWithSubscription))

	if err := CompareStructs(*userWithSubscriptionCheck, *userWithSubscription, []string{"TelegramID", "TelegramName", "CreatedAt", "UpdatedAt", "SubscriptionStatus", "SubscriptionExpiredAt"}); err != nil {
		t.Fatal(err)
	}

}

func hashStringToInt64(s string) int64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return int64(h.Sum64())
}
