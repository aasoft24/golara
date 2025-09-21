package observers

import (
	"fmt"

	"github.com/aasoft24/golara/app/models"
	"github.com/aasoft24/golara/wpkg/orm"
	"gorm.io/gorm"
)

type UserObserver struct{}

// User create হলে ইমেইল পাঠানো
func (UserObserver) Created(tx *gorm.DB, model interface{}) {
	user, ok := model.(*models.User)
	if !ok {
		return
	}
	fmt.Println("📧 Sending welcome email to:", user.Email)

}

// User update হলে লগ রাখা
func (UserObserver) Updated(tx *gorm.DB, model interface{}) {
	user, ok := model.(*models.User)
	if !ok {
		return
	}
	fmt.Println("📝 User updated:", user.ID, "Name:", user.Name)
	// এখানে তুমি চাইলে activity log table-এ insert করতে পারো
}

// User delete হলে কিছু করা
func (UserObserver) Deleted(tx *gorm.DB, model interface{}) {
	user, ok := model.(*models.User)
	if !ok {
		return
	}
	fmt.Println("❌ User deleted:", user.ID)
}

func init() {
	// Observer রেজিস্ট্রেশন
	orm.RegisterObserver("User", UserObserver{})
}
