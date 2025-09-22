package repository

type User struct {
	Id       int64  `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"size:50;unique;"`
	Password string `db:"password"`
	Mobile   string `gorm:"size:30;index:idx_name"`
}

type Order struct {
	Id       int64  `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"size:50;unique;"`
	Price    int64  `gorm:"index:idx_price"`
}
