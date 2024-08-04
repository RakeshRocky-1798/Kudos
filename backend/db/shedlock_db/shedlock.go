package shedlock_db

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type ShedlockRepository struct {
	logger     *zap.Logger
	gormClient *gorm.DB
}

func NewShedlockRepository(gormClient *gorm.DB, logger *zap.Logger) *ShedlockRepository {
	return &ShedlockRepository{
		gormClient: gormClient,
		logger:     logger,
	}
}

func (repo *ShedlockRepository) Insert(name string, lockedBy string, lockTime int) bool {

	var shedlockEntity ShedlockEntity
	s := &ShedlockEntity{
		Name:        name,
		LockUntil:   time.Now().Add(time.Duration(lockTime) * time.Second),
		LockedAt:    time.Now(),
		LockedBy:    lockedBy,
		LockedValue: true,
	}

	result := repo.gormClient.Table("shedlock").Where("name=?", name).Where("locked_value=?", false).Delete(&shedlockEntity)
	create := repo.gormClient.Create(&s)
	if create.Error == nil && create.RowsAffected > 0 {
		return true
	}

	if result.Error != nil {
		repo.logger.Error("Error While locking: ", zap.Error(result.Error))
	}
	return false
}

func (repo *ShedlockRepository) Unlock(name string, unlockedBy string) bool {
	now := time.Now()

	updatedData := map[string]interface{}{
		"locked_value": false,
		"lock_until":   now,
		"unlocked_by":  unlockedBy,
	}
	result := repo.gormClient.Table("shedlock").Where("name=?", name).Where("locked_value=?", true).UpdateColumns(updatedData)

	if result.Error != nil {
		repo.logger.Error("Failed to unlock ShedLock", zap.Error(result.Error))

		return false
	}

	if result.RowsAffected == 0 {
		repo.logger.Info("No ShedLock found with name: " + name + " for unlocking by: " + unlockedBy)
		return false
	}
	return true
}
