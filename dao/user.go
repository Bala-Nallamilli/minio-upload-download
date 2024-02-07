package dao

import "gorm.io/gorm"

type User struct {
	db *gorm.DB
}

func (u *User) Create(value User) error {
	return nil
}
func Update(value User) error {
	return nil
}

func Delete(value User) error {
	return nil
}

func Find(dest interface{}, conds ...interface{}) (User, error) {
	return User{}, nil
}
func FindAll(dest interface{}, conds ...interface{}) ([]User, error) {
	return nil, nil
}
