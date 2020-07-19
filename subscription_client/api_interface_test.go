package subscription_client

import (
	"fmt"
	"testing"
)

func TestSubscriptionClient_CreateDestination(t *testing.T) {
	s := NewSubscriptionClient("http://127.0.0.1:8000", "")
	resource, err := s.CreateDestination("12345", "discord")
	if err != nil {
		t.Error(err)
	}

	if resource.ID == 0 {
		t.Fail()
	}
}

func TestSubscriptionClient_CreateSubscriptionType(t *testing.T) {
	s := NewSubscriptionClient("http://127.0.0.1:8000", "")
	resource, err := s.CreateSubscriptionType("e621", "gay")
	if err != nil {
		t.Error(err)
	}

	if resource.ID == 0 {
		t.Fail()
	}
}

func TestSubscriptionClient_CreateSubscription(t *testing.T) {
	s := NewSubscriptionClient("http://127.0.0.1:8000", "")

	sub, err := s.CreateSubscription(21, 1)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(sub)
	if sub.ID == 0  {
		t.Fail()
	}
}

func TestSubscriptionClient_Subscribe(t *testing.T) {
	s := NewSubscriptionClient("http://127.0.0.1:8000", "")

	sub, err := s.Subscribe("discord", "1234", "e621", "gay")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(sub.SubscriptionType)
	if sub.ID == 0 {
		t.Fail()
	}
}

func TestSubscriptionClient_FindChannelSubscriptions(t *testing.T) {
	s := NewSubscriptionClient("http://127.0.0.1:8000", "")

	subs, err := s.FindChannelSubscriptions("1234")

	if err != nil {
		t.Error(err)
	}

	fmt.Println(len(subs))
}