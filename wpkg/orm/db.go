package orm

import (
	"github.com/aasoft24/golara/wpkg/database"
	"gorm.io/gorm"
)

type DBQuery struct {
	tx *gorm.DB
}

func DB() *DBQuery {
	return &DBQuery{tx: database.DB}
}

func (q *DBQuery) Table(name string) *DBQuery {
	q.tx = q.tx.Table(name)
	return q
}

func (q *DBQuery) Where(query interface{}, args ...interface{}) *DBQuery {
	q.tx = q.tx.Where(query, args...)
	return q
}

func (q *DBQuery) OrderBy(order string) *DBQuery {
	q.tx = q.tx.Order(order)
	return q
}

func (q *DBQuery) Limit(limit int) *DBQuery {
	q.tx = q.tx.Limit(limit)
	return q
}

func (q *DBQuery) Get() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := q.tx.Find(&results).Error
	return results, err
}

func (q *DBQuery) First() (map[string]interface{}, error) {
	var result map[string]interface{}
	err := q.tx.First(&result).Error
	return result, err
}

func (q *DBQuery) Insert(data map[string]interface{}) error {
	return q.tx.Create(data).Error
}

func (q *DBQuery) Update(data map[string]interface{}) error {
	return q.tx.Updates(data).Error
}

func (q *DBQuery) Delete() error {
	return q.tx.Delete(nil).Error
}

func (q *DBQuery) Raw(sql string, values ...interface{}) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := q.tx.Raw(sql, values...).Scan(&results).Error
	return results, err
}

// Transaction Methods
func (q *DBQuery) Begin() *DBQuery {
	q.tx = q.tx.Begin()
	return q
}

func (q *DBQuery) Commit() error {
	return q.tx.Commit().Error
}

func (q *DBQuery) Rollback() error {
	return q.tx.Rollback().Error
}

func (q *DBQuery) Macro(name string, args ...interface{}) *DBQuery {
	q.tx = CallMacro(q.tx, name, args...)
	return q
}
