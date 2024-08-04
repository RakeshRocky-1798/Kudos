package shedlock_db

import "time"

type ShedlockEntity struct {
	Name        string    `gorm:"name"`
	LockUntil   time.Time `gorm:"lock_until"`
	LockedAt    time.Time `gorm:"locked_at"`
	LockedBy    string    `gorm:"locked_by"`
	LockedValue bool      `gorm:"locked_value"`
	UnlockedBy  string    `gorm:"unlocked_by"`
}

func (ShedlockEntity) TableName() string {
	return "shedlock"
}
