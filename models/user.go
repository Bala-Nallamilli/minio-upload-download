package models

type User struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
}

func (User) TableName() string {
	return "users"
}
