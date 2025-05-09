package models

type Expression struct {
	ID         uint   `gorm:"primary_key"`
	Expression string `gorm:"not null"`
	Result     string `gorm:"default:'pending'"`
	UserID     uint
	User       User `gorm:"foreignkey:UserID"`
}
