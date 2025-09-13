package repository

type User struct {
	Id       int64  `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"size:50;unique;"`
}
