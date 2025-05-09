package models

type User struct {
	ID          uint         `gorm:"primary_key"`
	Username    string       `gorm:"unique;not null" json:"username"`
	Password    string       `gorm:"not null" json:"password"`
	Expressions []Expression `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
