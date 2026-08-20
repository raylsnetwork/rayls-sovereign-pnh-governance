package domain

import "gorm.io/gorm"

type PrivateNetwork struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
}

func (PrivateNetwork) TableName() string {
	return "private_networks"
}
