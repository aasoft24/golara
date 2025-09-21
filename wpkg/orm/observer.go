package orm

import "gorm.io/gorm"

type Observer interface {
	Created(tx *gorm.DB, model interface{})
	Updated(tx *gorm.DB, model interface{})
	Deleted(tx *gorm.DB, model interface{})
}

var observers = map[string][]Observer{}

func RegisterObserver(modelName string, obs Observer) {
	observers[modelName] = append(observers[modelName], obs)
}

func TriggerCreated(modelName string, tx *gorm.DB, model interface{}) {
	for _, obs := range observers[modelName] {
		obs.Created(tx, model)
	}
}

func TriggerUpdated(modelName string, tx *gorm.DB, model interface{}) {
	for _, obs := range observers[modelName] {
		obs.Updated(tx, model)
	}
}

func TriggerDeleted(modelName string, tx *gorm.DB, model interface{}) {
	for _, obs := range observers[modelName] {
		obs.Deleted(tx, model)
	}
}
