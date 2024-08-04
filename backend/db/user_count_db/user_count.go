package userCountDb

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserCountRepository struct {
	logger     *zap.Logger
	gormClient *gorm.DB
}

func NewUserCountRepository(logger *zap.Logger, gormClient *gorm.DB) *UserCountRepository {
	return &UserCountRepository{
		logger:     logger,
		gormClient: gormClient,
	}
}

func (r *UserCountRepository) AddUserCountEntity(request *CreateUserCountRequest) (*UserCountEntity, error) {

	userCountEntity := &UserCountEntity{
		UserId:        request.UserId,
		GivenCount:    request.GivenCount,
		ReceivedCount: request.ReceivedCount,
		Month:         request.Month,
		Week:          request.Week,
	}

	result := r.gormClient.Create(userCountEntity)
	if result.Error != nil {
		return nil, result.Error
	}

	return userCountEntity, nil
}

func (r *UserCountRepository) CurrentWeekCount(userID string, currentWeek string) (*UserCountEntity, error) {

	if r == nil || r.gormClient == nil {
		r.logger.Error("[CurrentWeekCount]Repository or gormClient is nil")
		return nil, nil
	}
	var userCountGroup UserCountEntity
	result := r.gormClient.Table("user_count").Select("*").
		Where("user_id = ? AND week = ?", userID, currentWeek).
		Scan(&userCountGroup)

	if result.Error != nil {
		return nil, result.Error
	}
	return &userCountGroup, nil
}

func (r *UserCountRepository) CurrentMonthCount(userID string, currentMonth string) (*UserMonthlyCount, error) {

	if r == nil || r.gormClient == nil {
		r.logger.Error("[CurrentMonthCount]Repository or gormClient is nil")
		return nil, nil
	}
	var userCountGroup UserMonthlyCount
	result := r.gormClient.Table("user_count").Select("user_id, count(*) as count__").
		Where("user_id = ? AND month = ?", userID, currentMonth).
		Group("user_id").Scan(&userCountGroup)

	if result.Error != nil {
		return nil, result.Error
	}
	return &userCountGroup, nil
}

func (r *UserCountRepository) UpdateGivenCount(userId string, week string) {
	if r == nil || r.gormClient == nil {
		r.logger.Error("[updateGivenCount]Repository or gormClient is nil")
		return
	}
	sql := `UPDATE user_count SET given_count = given_count + 1, updated_at= now() WHERE user_id = ? AND week = ?`
	result := r.gormClient.Exec(sql, userId, week)
	if result.Error != nil {
		r.logger.Error("Error while setting the given count", zap.Error(result.Error))
		return
	}
}

func (r *UserCountRepository) UpdateReceivedCount(userId string, week string) {
	if r == nil || r.gormClient == nil {
		r.logger.Error("[updateGivenCount]Repository or gormClient is nil")
		return
	}
	sql := `UPDATE user_count SET received_count = received_count + 1, updated_at= now() WHERE user_id = ? AND week = ?`
	result := r.gormClient.Exec(sql, userId, week)
	if result.Error != nil {
		r.logger.Error("Error while setting the received count", zap.Error(result.Error))
		return
	}
}

func (r *UserCountRepository) GetUserGiveCountCurrentWeek(userId string, week string) (*UserGivenCount, error) {
	if r == nil || r.gormClient == nil {
		r.logger.Error("[updateGivenCount]Repository or gormClient is nil")
		return nil, nil
	}

	var userCountGroup UserGivenCount
	result := r.gormClient.Table("user_count").
		Select("user_id, given_count").
		Where("user_id=? and week=?", userId, week).
		Scan(&userCountGroup)

	if result.Error != nil {
		r.logger.Error("Error while setting the week count", zap.Error(result.Error))
		return nil, result.Error
	}
	return &userCountGroup, nil
}

func (r *UserCountRepository) GetUserLeaderBoardData(isReceived bool, month string) (*[]LeaderBoardData, error) {

	var LeaderBoardGroup []LeaderBoardData
	var result *gorm.DB

	if isReceived {
		result = r.gormClient.Table("user_count").
			Select("user_id, SUM(received_count) as count__, ROW_NUMBER() OVER (ORDER BY SUM(received_count) DESC) AS rank").
			Where("month = ?", month).
			Group("user_id").
			Having("SUM(received_count) > 0").
			Order("SUM(received_count) DESC").
			Scan(&LeaderBoardGroup)
	} else {
		result = r.gormClient.Table("user_count").
			Select("user_id, SUM(given_count) as count__, ROW_NUMBER() OVER (ORDER BY SUM(given_count) DESC) AS rank").
			Where("month = ?", month).
			Group("user_id").
			Having("SUM(given_count) > 0").
			Order("SUM(given_count) DESC").
			Scan(&LeaderBoardGroup)
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &LeaderBoardGroup, nil
}
