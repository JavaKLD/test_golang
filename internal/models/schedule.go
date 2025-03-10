package models

import "time"

type Schedule struct {//график
	ID uint `gorm:"primaryKey;autoIncrement" json:"schedule_id"`
	Aid_name string `gorm:"not null" json:"aid_name"`
	Aid_per_day int `gorm:"not null" json:"aid_per_day"`
	UserID uint `gorm:"not null" json:"user_id"`
	Duration int `json:"duration"`
	Create_at time.Time `json:"create_at" gorm:"autoCreateTime"`
}
