package example

import (
	"gorm.io/gorm"
)

type ExampleRepo struct {
	Db *gorm.DB
}

func NewExampleRepo() *ExampleRepo {
	return &ExampleRepo{}
}

func (r *ExampleRepo) CreateExample(s interface{}) error {
	tx := r.Db.Create(s)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *ExampleRepo) UpdateExample() {
	return
}

func (r *ExampleRepo) GetExample() {
	return
}

func (r *ExampleRepo) DeleteExample() {
	return
}
