package models

import "time"

type Schedule struct {
	CreatedAt time.Time `json:"create_at" gorm:"type:datetime(0);autoCreateTime"`
	AidName   string    `gorm:"not null; " json:"aid_name"`
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"schedule_id"`
	AidPerDay uint64    `gorm:"not null" json:"aid_per_day"`
	UserID    uint64    `gorm:"not null" json:"user_id"`
	Duration  int64     `json:"duration"`
}
