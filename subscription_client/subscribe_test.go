package subscription_client

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/plally/FoxBot/subscription_client/subtypes"
	_ "github.com/plally/FoxBot/subscription_client/subtypes"
	"github.com/plally/subscription_api/database"
	"github.com/sirupsen/logrus"
	"log"
	"testing"
)

func makesubclient() (*SubscriptionClient){
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"127.0.0.1",
		"5432",
		"dev",
		"fox",
		"fox_bot_dev",
	)

	db, err :=gorm.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	logrus.SetLevel(logrus.DebugLevel)
	db = db.LogMode(false)

	database.Migrate(db)

	subClient := NewSubscriptionClient(db)
	subtypes.RegisterE621()
	subtypes.RegisterRSS()
	return subClient
}

func teardown(subClient *SubscriptionClient) {
	subClient.DB.Exec("DROP TABLE destinations CASCADE;")
	subClient.DB.Exec("DROP TABLE subscriptions CASCADE;")
	subClient.DB.Exec("DROP TABLE subscription_types CASCADE;")
}

func TestSubscriptionClient_Subscribe(t *testing.T) {
	subClient := makesubclient()

	testData := []struct {
		subType string
		tags    string
		dest    string
	}{
	{"rss", "https://localhost:2134/feed.rss", "1234"},
	{"rss", "https://localhost:2134/feed.rss", "12432"},
	}

	uniqueSubs := make(map[uint]*database.Subscription)
	for _, data := range testData {
		sub, err := subClient.Subscribe(data.subType, data.tags, data.dest)
		if err != nil { t.Error(err) }
		uniqueSubs[sub.ID] = sub
		if sub.ID <= 0 {
			t.Fail()
		}
	}

	if len(uniqueSubs) != len(testData) {
		t.Fail()
	}

	teardown(subClient)
}


func TestSubscriptionClient_GetSubscription(t *testing.T) {
	subClient := makesubclient()

	testData := []struct {
		subType string
		tags    string
		dest    string
	}{
		{"rss", "https://housepetscomic.com/feed", "123"},
		{"rss", "http://example.com/feed.rss", "123"},
	}

	for _, data := range testData {
		sub, err:= subClient.Subscribe(data.subType, data.tags, data.dest)
		if err != nil { t.Error(err) }
		if sub.ID <= 0 {
			t.Fail()
		}
	}
	subs, err := subClient.GetSubscriptions("123")
	if err != nil { t.Error(err) }
	if len(subs) < len(testData) {
		t.Fail()
	}

	teardown(subClient)
}

func TestSubscriptionClient_DeleteSubscription(t *testing.T) {
	subClient := makesubclient()

	// create a subscription
	sub, err:= subClient.Subscribe("rss", "http://example.com/feed.rss", "312")
	if err != nil { t.Error(err) }

	if sub.ID <= 0 {
		t.Fail()
	}

	// delete
	err  = subClient.DeleteSubscription("rss", "http://example.com/feed.rss", "312")
	if err != nil {
		t.Error(err)
	}
	// there should now be no more subs with that destination left
	subs, err := subClient.GetSubscriptions("312")
	if err != nil {
		t.Error(err)
	}
	if len(subs) != 0 {
		t.Fail()
	}

	teardown(subClient)
}