package achiementDb

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AchievementRepository struct {
	logger     *zap.Logger
	gormClient *gorm.DB
}

func NewAchievementRepository(logger *zap.Logger, gormClient *gorm.DB) *AchievementRepository {
	return &AchievementRepository{
		logger:     logger,
		gormClient: gormClient,
	}
}

func (a *AchievementRepository) GetAllAchievementName() ([]AchievementEntity, error) {

	if a == nil || a.gormClient == nil {
		a.logger.Error("[Achievement]Repository or gormClient is nil")
		return nil, errors.New("[Achievement]repository or gormClient is nil")
	}

	var achievementEntity []AchievementEntity

	var query = a.gormClient.Model([]AchievementEntity{})
	query = query.Order("id asc").Find(&achievementEntity)
	if query.Error != nil {
		a.logger.Error("[Achievement]Error while fetching achievement data", zap.Error(query.Error))
		return nil, query.Error
	}
	return achievementEntity, nil
}

func (a *AchievementRepository) GetAchievementIdFromName(name string) (*AchievementName, error) {

	if a == nil || a.gormClient == nil {
		a.logger.Error("[Achievement]Repository or gormClient is nil")
		return nil, errors.New("[Achievement]repository or gormClient is nil")
	}

	var achievementEntity AchievementName
	result := a.gormClient.Table("achievement").Select("id, achievement_name").
		Where("achievement_name=?", name).Scan(&achievementEntity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &achievementEntity, nil
}

func (a *AchievementRepository) GetAchievementNameFromID(id int) (*AchievementName, error) {

	if a == nil || a.gormClient == nil {
		a.logger.Error("[Achievement]Repository or gormClient is nil")
		return nil, errors.New("[Achievement]repository or gormClient is nil")
	}

	var achievementEntity AchievementName
	result := a.gormClient.Table("achievement").Select("id, achievement_name").
		Where("id=?", id).Scan(&achievementEntity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &achievementEntity, nil
}

func (a *AchievementRepository) CheckIfAchievementExists(id int) (bool, error) {

	if a == nil || a.gormClient == nil {
		a.logger.Error("[Achievement]Repository or gormClient is nil")
		return false, errors.New("[Achievement]repository or gormClient is nil")
	}

	var count int64
	result := a.gormClient.Table("achievement").Where("id=?", id).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
