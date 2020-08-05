package subscription_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/plally/subscription_api/database"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type SubscriptionClient struct {
	BaseURL string
	Client  *http.Client
	Token   string
}

func NewSubscriptionClient(baseURL string, token string) *SubscriptionClient {
	return &SubscriptionClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
		Token:   token,
	}
}

func (s *SubscriptionClient) Request(method string, endpoint string, body io.Reader, params map[string]string) (*http.Response, error) {
	url := s.BaseURL + endpoint
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", s.Token)
	log.Debugf("Subscription Client Request made: %s %s", req.Method, req.URL.Path)
	resp, err := s.Client.Do(req)

	return resp, err
}

func (s *SubscriptionClient) PostJson(endpoint string, object interface{}) (*http.Response, error) {
	data, _ := json.Marshal(object)
	body := bytes.NewReader(data)
	return s.Request("POST", endpoint, body, map[string]string{})
}

func (s *SubscriptionClient) Create(endpoint string, outputModel interface{}, inputData interface{}) error {
	resp, err := s.PostJson(endpoint, inputData)

	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case http.StatusForbidden:
		return ErrNoPermissions
	case http.StatusConflict:
		return ErrAlreadyExists
	case http.StatusCreated:
		fallthrough
	case http.StatusOK:
		break
	default:
		return SubError{fmt.Sprintf("Unknown error creating resource: %v", resp.StatusCode), nil}
	}

	data, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, outputModel)

	return nil
}

func (s *SubscriptionClient) CreateDestination(ident, destType string) (database.Destination, error) {
	var dest database.Destination

	err := s.Create("/destinations", &dest, map[string]interface{}{
		"external_identifier": ident,
		"destination_type":    destType,
	})

	return dest, err
}

func (s *SubscriptionClient) CreateSubscriptionType(typeString, tags string) (database.SubscriptionType, error) {
	var subType database.SubscriptionType

	err := s.Create("/subscription_types", &subType, map[string]interface{}{
		"type": typeString,
		"tags": tags,
	})

	return subType, err
}

func (s *SubscriptionClient) CreateSubscription(destination, subtype uint) (database.Subscription, error) {
	var sub database.Subscription

	err := s.Create("/subscriptions", &sub, map[string]uint{
		"destination_id":       destination,
		"subscription_type_id": subtype,
	})

	return sub, err
}

func (s *SubscriptionClient) Subscribe(destType, destIdent, subType, subTags string) (database.Subscription, error) {
	var sub database.Subscription

	err := s.Create("/subscribe", &sub, map[string]string{
		"destination_type":       destType,
		"destination_identifier": destIdent,
		"subscription_type":      subType,
		"subscription_tags":      subTags,
	})

	return sub, err
}

func (s *SubscriptionClient) FindChannelSubscriptions(channelid string) ([]database.Subscription, error) {
	var destinations []database.Destination
	var subscriptions []database.Subscription

	err := s.Find("/destinations", &destinations, map[string]string{
		"external_identifier": channelid,
		"destination_type":    "discord",
	})

	if err != nil || len(destinations) == 0 {
		return subscriptions, err
	}

	dest := destinations[0]

	err = s.Find("/subscriptions", &subscriptions, map[string]string{
		"destination_id": strconv.FormatUint(uint64(dest.ID), 10),
	})

	return subscriptions, err
}

func (s *SubscriptionClient) Find(endpoint string, model interface{}, params map[string]string) error {
	resp, err := s.Request("GET", endpoint, nil, params)
	if err != nil {
		return err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	err = json.Unmarshal(data, model)
	return err
}

func (s *SubscriptionClient) DeleteSubscription(id int) (database.Subscription, error) {
	var sub database.Subscription
	err := s.Delete("/subscriptions", &sub, id)
	return sub, err

}
func (s *SubscriptionClient) Delete(endpoint string, model interface{}, id int) error {
	resp, err := s.Request("DELETE", fmt.Sprintf("%v/%v", endpoint, id), nil, map[string]string{})
	if err != nil {
		return err
	}

	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	err = json.Unmarshal(data, model)
	return err
}
