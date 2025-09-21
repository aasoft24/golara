package models

import (
	"fmt"
	"time"

	"github.com/aasoft24/golara/wpkg/orm"
	"gorm.io/gorm"
)

type User struct {
	ID        int64 `gorm:"primaryKey;autoIncrement;column:id"`
	Name      string
	Mobile    string
	Email     string
	Password  string    `json:"--"`
	Balance   float64   `gorm:"type:decimal(10,2);default:0"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `db:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt time.Time `db:"deleted_at" gorm:"autoDeleteTime"`
}

func init() {
	orm.RegisterModel(User{})
}

func UserModel() *orm.Query[User] {
	return orm.NewQuery(User{})
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("Before creating user:", u.Name)
	return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Println("After creating user:", u.ID)
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	fmt.Println("Before updating user:", u.ID)
	return
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	fmt.Println("After updating user:", u.ID)
	return
}

func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Println("Before deleting user:", u.ID)
	return
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Println("After deleting user:", u.ID)
	return
}
