package entities

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Visit struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Visit_uid   string         `gorm:"index;type:varchar(22);primaryKey"`
	Doctor_uid  string         `gorm:"index;type:varchar(22)"`
	Clinic_uid  string         `gorm:"index;type:varchar(22)"`
	Patient_uid string         `gorm:"index;type:varchar(22)"`
	Date        datatypes.Date
	Status      string    `gorm:"type:enum('waiting', 'done', 'cancel');default:'waiting'"`
}