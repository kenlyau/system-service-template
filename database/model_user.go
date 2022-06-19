package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `gorm:"unique"`
	Age   int
	Phone string
	Addr  string
}

type Users []User

func ModelUserAdd(u *User) error {
	return DB.Create(u).Error
}

func ModelUsersQuery() (Users, error) {
	var users Users
	err := DB.Find(&users).Error
	return users, err
}
