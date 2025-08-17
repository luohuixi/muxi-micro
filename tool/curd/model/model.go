package model

type User struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	Username  string `gorm:"size:50;primaryKey;"`
	FirstName string `gorm:"size:30;index:idx_name"`
}
