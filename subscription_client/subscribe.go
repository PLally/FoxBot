package subscription_client

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/plally/FoxBot/subscription_client/subtypes"
	"github.com/plally/subscription_api/database"
	"github.com/plally/subscription_api/subscription"
)

type SubscriptionClient struct {
	DB *gorm.DB
}

func NewSubscriptionClient(db *gorm.DB) (*SubscriptionClient) {
	database.Migrate(db)

	return &SubscriptionClient{
		DB: db,
	}
}

func (s *SubscriptionClient) Subscribe(subType string, tags string, dest string) (*database.Subscription, error) {
	handler := subscription.GetSubTypeHandler(subType)
	if handler == nil {
		return nil, SubError{"Invalid subscription type", nil }
	}
	tags, err := handler.Validate(tags)
	if err != nil {
		return nil, SubError{ "Invalid subscription tags", err }
	}

	subscriptionType := database.SubscriptionType{
		Type: subType,
		Tags: tags,
	}

	destination := database.Destination{
		ExternalIdentifier:dest,
		DestinationType:"discord",
	}

	if err := s.DB.FirstOrCreate(&subscriptionType, subscriptionType).Error; err != nil {
		return nil, err
	}
	if err := s.DB.FirstOrCreate(&destination, destination).Error; err != nil {
		return nil, err
	}

	sub := database.Subscription{
		DestinationID: destination.ID,
		SubscriptionTypeID: subscriptionType.ID,
	}

	s.DB.First(&sub, sub)
	if sub.ID != 0 {
		return &sub, SubError{fmt.Sprintf("Subscription %v : %v already exists", subType, tags), nil }
	}
	err = s.DB.Create(&sub).Error
	return &sub, err
}

func (s *SubscriptionClient) GetSubscriptions(dest string) ([]database.Subscription, error){
	var subs []database.Subscription

	destination := &database.Destination{}

	err := s.DB.Where("external_identifier=?", dest).Find(destination).Error

	if err != nil { return nil, err }

	s.DB.Preload("Destination").
		Preload("SubscriptionType").
		Where("destination_id=?", destination.ID).
		Find(&subs)

	return subs, nil
}

func (s *SubscriptionClient) DeleteSubscription(subtype string, tags string, dest string) error {
	handler := subscription.GetSubTypeHandler(subtype)
	if handler == nil { return SubError{"Invalid sub type", nil}}
	tags, err := handler.Validate(tags)
	if err != nil { return err }

	destination := &database.Destination{}
	if err := s.DB.Where("external_identifier=?", dest).Find(destination).Error; err != nil {
		return SubError{"Couldn't find the destination "+dest, err }
	}

	subscriptionType := database.SubscriptionType{}
	if err := s.DB.Where("type=? AND tags=?", subtype, tags).Find(&subscriptionType).Error; err != nil {
		return SubError{fmt.Sprintf("Couldnt find subscription type %v : %v", subtype, tags), err }
	}

	sub := database.Subscription{
		SubscriptionTypeID: subscriptionType.ID,
		DestinationID: destination.ID,
	}
	err = s.DB.Delete(&sub).Error

	if err != nil {
		return SubError{ "Could not delete the subscription", err }
	}
	return nil
}

type SubError struct {
	s string // a readable message to show to the user
	err error
}

func (e SubError) Error() string {
	return e.s
}

func (e SubError) Unwrap() error {
	return e.err
}