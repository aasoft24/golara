package orm

import "gorm.io/gorm"

type Scope func(*gorm.DB) *gorm.DB

func ApplyScopes(db *gorm.DB, scopes ...Scope) *gorm.DB {
	for _, scope := range scopes {
		db = scope(db)
	}
	return db
}
