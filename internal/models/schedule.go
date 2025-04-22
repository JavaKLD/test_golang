package models

import "time"

type Schedule struct { //график
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"schedule_id"`
	Aid_name    string    `gorm:"not null; " json:"aid_name"`
	Aid_per_day uint64    `gorm:"not null" json:"aid_per_day"`
	UserID      uint64    `gorm:"not null" json:"user_id"`
	Duration    int64     `json:"duration"`
	Create_at   time.Time `json:"create_at" gorm:"type:datetime(0);autoCreateTime"`
}
