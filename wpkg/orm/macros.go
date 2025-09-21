package orm

import "gorm.io/gorm"

type Macro func(*gorm.DB, ...interface{}) *gorm.DB

var macros = map[string]Macro{}

func RegisterMacro(name string, fn Macro) {
	macros[name] = fn
}

func CallMacro(db *gorm.DB, name string, args ...interface{}) *gorm.DB {
	if macro, ok := macros[name]; ok {
		return macro(db, args...)
	}
	return db
}
