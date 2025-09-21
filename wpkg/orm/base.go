package orm

import (
	"github.com/aasoft24/golara/wpkg/database"
	"gorm.io/gorm"
)

type Query[T any] struct {
	db *gorm.DB
}

func NewQuery[T any](model T) *Query[T] {
	return &Query[T]{db: database.DB.Model(&model)}
}

func (q *Query[T]) Where(query interface{}, args ...interface{}) *Query[T] {
	q.db = q.db.Where(query, args...)
	return q
}

func (q *Query[T]) OrderBy(order string) *Query[T] {
	q.db = q.db.Order(order)
	return q
}

func (q *Query[T]) Limit(limit int) *Query[T] {
	q.db = q.db.Limit(limit)
	return q
}

func (q *Query[T]) With(relations ...string) *Query[T] {
	for _, r := range relations {
		q.db = q.db.Preload(r)
	}
	return q
}

func (q *Query[T]) Get() ([]T, error) {
	var results []T
	err := q.db.Find(&results).Error
	return results, err
}

func (q *Query[T]) First() (T, error) {
	var result T
	err := q.db.First(&result).Error
	return result, err
}

func (q *Query[T]) Create(data *T) error {
	return q.db.Create(data).Error
}

func (q *Query[T]) Update(data interface{}) error {
	return q.db.Updates(data).Error
}

func (q *Query[T]) Delete() error {
	return q.db.Delete(nil).Error
}

// Shortcuts
func (q *Query[T]) All() ([]T, error) {
	var results []T
	err := q.db.Find(&results).Error
	return results, err
}

func (q *Query[T]) Find(id interface{}) (T, error) {
	var result T
	err := q.db.First(&result, id).Error
	return result, err
}

// Pagination
func (q *Query[T]) Paginate(page, pageSize int) ([]T, int64, error) {
	var results []T
	var total int64
	q.db.Count(&total)
	offset := (page - 1) * pageSize
	err := q.db.Offset(offset).Limit(pageSize).Find(&results).Error
	return results, total, err
}

// Relationship helpers
func (q *Query[T]) HasMany(model interface{}, foreignKey string) *gorm.DB {
	return q.db.Model(model).Association(foreignKey).DB
}

func (q *Query[T]) HasOne(model interface{}, foreignKey string) *gorm.DB {

	return q.db.Model(model).Association(foreignKey).DB
}

func (q *Query[T]) BelongsTo(model interface{}, foreignKey string) *gorm.DB {
	return q.db.Model(model).Association(foreignKey).DB
}

// Transaction Methods
func (q *Query[T]) Begin() *Query[T] {
	q.db = q.db.Begin()
	return q
}

func (q *Query[T]) Commit() error {
	return q.db.Commit().Error
}

func (q *Query[T]) Rollback() error {
	return q.db.Rollback().Error
}
func (q *Query[T]) Macro(name string, args ...interface{}) *Query[T] {
	q.db = CallMacro(q.db, name, args...)
	return q
}
