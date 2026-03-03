package postgres

import (
	"time"

	"gorm.io/gorm"
)

type UserRecord struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Username  string    `gorm:"column:username;uniqueIndex;not null;type:varchar(255)"`
	Password  string    `gorm:"column:password;not null;type:varchar(255)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeleteAt  gorm.DeletedAt
}

func (UserRecord) TableName() string {
	return "users"
}
