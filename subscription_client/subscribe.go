package subscription_client

import (
	"errors"
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
var (
	InvalidSubType = errors.New("Invalid subscription type")
)
func (s *SubscriptionClient) Subscribe(subType string, tags string, dest string) (*database.Subscription, error) {
	handler := subscription.GetSubTypeHandler(subType)
	if handler == nil { return nil,  InvalidSubType }
	tags, err := handler.Validate(tags)
	if err != nil { return nil, err }

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

	err = s.DB.FirstOrCreate(&sub, sub).Error
	return &sub, err
}

func (s *SubscriptionClient) GetSubscriptions(dest string) ([]database.Subscription, error){
	var subs []database.Subscription
	destination := &database.Destination{}
	err := s.DB.Where("external_identifier=?", dest).Find(destination).Error
	if err != nil { return nil, err }
	err = s.DB.Preload("Destination").
		Preload("SubscriptionType").
		Where("destination_id=?", destination.ID).
		Find(&subs).Error

	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (s *SubscriptionClient) DeleteSubscription(subtype string, tags string, dest string) error {
	handler := subscription.GetSubTypeHandler(subtype)
	if handler == nil { return InvalidSubType }
	tags, err := handler.Validate(tags)
	if err != nil { return err }

	destination := &database.Destination{}
	if err := s.DB.Where("external_identifier=?", dest).Find(destination).Error; err != nil {
		return err
	}

	subscriptionType := database.SubscriptionType{}
	if err := s.DB.Where("type=? AND tags=?", subtype, tags).Find(&subscriptionType).Error; err != nil {
		return err
	}

	subscription := database.Subscription{
		SubscriptionTypeID: subscriptionType.ID,
		DestinationID: destination.ID,
	}
	err = s.DB.Delete(&subscription).Error
	return err
}
