package entities

import (
	"time"

	"gorm.io/gorm"
)

type Doctor struct {
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	Doctor_uid     string         `gorm:"index;type:varchar(22);primaryKey"`
	Doctor_uid_ref string         `gorm:"index;type:varchar(22)"`
	UserName       string         `gorm:"index;not null;type:varchar(100)"`
	Email          string         `gorm:"index;not null;type:varchar(100)"`
	Password       string         `gorm:"not null;type:varchar(100)"`
	Name           string
	Image          string `gorm:"default:'https://www.teralogistics.com/wp-content/uploads/2020/12/default.png'"`
	Address        string
	Status         string `gorm:"type:enum('available', 'unAvailable');default:'available'"`
	OpenDay        string `gorm:"type:enum('senin', 'selasa', 'rabu', 'kamis', 'jumat', 'sabtu', 'minggu');default:'senin'"`
	CloseDay       string `gorm:"type:enum('senin', 'selasa', 'rabu', 'kamis', 'jumat', 'sabtu', 'minggu');default:'senin'"`
	Capacity       int
	Type           string  `gorm:"type:enum('doctor', 'admin');default:'doctor'"`
	Visits         []Visit `gorm:"foreignKey:Doctor_uid;references:Doctor_uid"`
}
