package kleosDb

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type KleosRepository struct {
	logger     *zap.Logger
	gormClient *gorm.DB
}

func NewKleosRepository(logger *zap.Logger, gormClient *gorm.DB) *KleosRepository {
	return &KleosRepository{
		logger:     logger,
		gormClient: gormClient,
	}
}

func (r *KleosRepository) GiveKleos(request *CreateKleosRequest) (*MyKleosEntity, error) {

	kleosEntity := &MyKleosEntity{
		SenderID:    request.SenderID,
		Message:     request.Message,
		Achievement: request.Achievement,
		ReceiverID:  request.ReceiverID,
		Year:        request.Year,
		Month:       request.Month,
		Week:        request.Week,
		Day:         request.Day,
	}

	result := r.gormClient.Create(kleosEntity)
	if result.Error != nil {
		return nil, result.Error
	}

	return kleosEntity, nil
}

func (r *KleosRepository) GetAllKleosData() (*[]MyKleosEntity, error) {
	var kleosGroup []MyKleosEntity
	result := r.gormClient.Table("kleos").Select("*").
		Order("created_at desc").Scan(&kleosGroup)
	if result.Error != nil {
		return nil, result.Error
	}
	return &kleosGroup, nil

}

func (r *KleosRepository) KleosGiverByUserCurrentWeek(userID string, currentWeek string) (*[]KleosGiven, error) {
	var kleosGroup []KleosGiven
	result := r.gormClient.Table("kleos").Select("sender_id, count(*) as count__").
		Where("sender_id = ? AND week = ?", userID, currentWeek).
		Group("sender_id").Scan(&kleosGroup)
	if result.Error != nil {
		return nil, result.Error
	}
	return &kleosGroup, nil
}

func (r *KleosRepository) GetMyKleosPerAchievement(userID string) (*[]KleosAchievementStruct, error) {
	var kleosGroup []KleosAchievementStruct
	result := r.gormClient.Table("kleos").Select("achievement, count(*) as count__").
		Where("receiver_id = ?", userID).
		Group("achievement").Order("count(*) DESC").Scan(&kleosGroup)
	if result.Error != nil {
		return nil, result.Error
	}
	return &kleosGroup, nil
}

func (r *KleosRepository) KleosGivenCountPerUser(userID string) (*KleosGiven, error) {
	var kleosGroup KleosGiven
	result := r.gormClient.Table("kleos").Select("sender_id, count(*) as count__").
		Where("sender_id = ?", userID).
		Group("sender_id").Scan(&kleosGroup)
	if result.Error != nil {
		return nil, result.Error
	}
	return &kleosGroup, nil
}

func (r *KleosRepository) KleosReceivedPerUser(userID string) (*KleosReceived, error) {
	var kleosGroup KleosReceived
	result := r.gormClient.Table("kleos").Select("receiver_id, count(*) as count__").
		Where("receiver_id = ?", userID).
		Group("receiver_id").Scan(&kleosGroup)
	if result.Error != nil {
		return nil, result.Error
	}
	return &kleosGroup, nil
}

func (r *KleosRepository) GetLeaderBoardBySlackIdUser(userID int) (*[]LeaderBoardGroupStruct, error) {
	var LeaderBoardGroup []LeaderBoardGroupStruct
	result := r.gormClient.Table("kleos").Select("receiver_id, count(*) as count__").
		Group("receiver_id").Order("count__ DESC").
		Limit(10).Scan(&LeaderBoardGroup)
	if result.Error != nil {
		return nil, result.Error
	}
	return &LeaderBoardGroup, nil
}

func (r *KleosRepository) LastThreeKleosData(userID string) (*[]LastThreeKleos, error) {
	var lastThreeKleos []LastThreeKleos
	result := r.gormClient.Table("kleos").Select("sender_id, message, achievement").
		Where("receiver_id = ?", userID).
		Order("id desc").
		Limit(3).Scan(&lastThreeKleos)
	if result.Error != nil {
		return nil, result.Error
	}
	return &lastThreeKleos, nil

}

func (r *KleosRepository) FetchPaginatedKleos(dataType string, pageNumber int64, pageSize int64, userId string) ([]FilteredKleosEntity, int, error) {

	var query = r.gormClient.Model([]FilteredKleosEntity{})
	var filteredKleosEntity []FilteredKleosEntity
	var totalElements int64

	switch dataType {
	case "given":
		{
			result := query.Where("sender_id = ?", userId).Count(&totalElements)
			if result.Error != nil {
				return nil, 0, result.Error
			}
			result = query.Select("id, sender_id, receiver_id, Achievement, Message, Created_at").
				Where("sender_id = ?", userId).Order("id desc").
				Offset(int(pageNumber * pageSize)).
				Limit(int(pageSize)).
				Find(&filteredKleosEntity)
			if result.Error != nil {
				return nil, 0, result.Error
			}
			if result.RowsAffected == 0 {
				return nil, int(totalElements), nil
			}
		}
	case "received":
		{
			result := query.Where("receiver_id = ?", userId).Count(&totalElements)
			if result.Error != nil {
				return nil, 0, result.Error
			}
			result = query.Select("id, sender_id, receiver_id, Achievement, Message, Created_at").
				Where("receiver_id = ?", userId).
				Order("id desc").
				Offset(int(pageNumber * pageSize)).
				Limit(int(pageSize)).Find(&filteredKleosEntity)

			if result.Error != nil {
				return nil, 0, result.Error
			}
			if result.RowsAffected == 0 {
				return nil, int(totalElements), nil
			}
		}
	default:
		{
			return nil, 0, errors.New("not a valid data type")
		}

	}
	return filteredKleosEntity, int(totalElements), nil
}

func (r *KleosRepository) GetLeaderBoardData(isReceived bool) (*[]LeaderBoardData, error) {
	var LeaderBoardGroup []LeaderBoardData
	var result *gorm.DB

	if isReceived {
		result = r.gormClient.Table("users").
			Select("id as user_id, received_count as count__, ROW_NUMBER() OVER (ORDER BY received_count DESC) AS rank").
			Order("count__ DESC").
			Scan(&LeaderBoardGroup)
	} else {
		result = r.gormClient.Table("users").
			Select("id as user_id, given_count as count__, ROW_NUMBER() OVER (ORDER BY given_count DESC) AS rank").
			Order("count__ DESC").
			Scan(&LeaderBoardGroup)
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &LeaderBoardGroup, nil
}

func (r *KleosRepository) GetCurrentUserData(isReceived bool, userId string) (*LeaderBoardData, error) {
	var currentUserData LeaderBoardData
	var result *gorm.DB

	if isReceived {
		result = r.gormClient.Table("kleos").
			Select("receiver_id as user_id, COUNT(*) as count__, ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) AS row_number").
			Group("receiver_id").Having("receiver_id = ?", userId).Order("count__ DESC").Limit(10).Scan(&currentUserData)
	} else {
		result = r.gormClient.Table("kleos").
			Select("sender_id as user_id, COUNT(*) as count__, ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) AS row_number").
			Group("sender_id").Having("sender_id = ?", userId).Order("count__ DESC").Limit(10).Scan(&currentUserData)
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &currentUserData, nil
}

func (r *KleosRepository) GetAdminData(dataType string) (*[]AdminAllUser, error) {

	var adminAllUser []AdminAllUser
	result := r.gormClient.Table("kleos")
	if dataType == "given" {
		result = result.Select("sender_id as user_id, COUNT(*) as count__").
			Group("sender_id").Order("count__ DESC").Scan(&adminAllUser)
	} else if dataType == "received" {
		result = result.Select("receiver_id as user_id, COUNT(*) as count__").
			Group("receiver_id").Order("count__ DESC").Scan(&adminAllUser)
	} else {
		return nil, errors.New("not a valid data type")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &adminAllUser, nil
}

func (r *KleosRepository) GetAdminDataTest() (*[]AdminAllUser, error) {

	var adminAllUser []AdminAllUser
	result := r.gormClient.Raw("SELECT user_id, kleos_given, kleos_received FROM " +
		"(select user_id, sum(given_kleos) as kleos_given, sum(received_kleos) as kleos_received FROM " +
		"(select sender_id as user_id, count(*) as given_kleos, 0 as received_kleos from kleos group by sender_id  UNION ALL " +
		"select receiver_id as user_id, 0 as given_kleos, count(*) as received_kleos from kleos group by receiver_id) as t1 " +
		"group by user_id) t2 " +
		"order by kleos_given desc").Scan(&adminAllUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &adminAllUser, nil
}

func (r *KleosRepository) GetPaginatedAdminData(pageNumber int64, pageSize int64) ([]AdminAllUser, int, error) {

	query := r.gormClient.Model([]AdminAllUser{})
	var adminAllUser []AdminAllUser
	var totalElements int64
	result := query.Count(&totalElements)

	result = r.gormClient.Raw("select user_id, sum(given_kleos) as kleos_given, sum(received_kleos) as kleos_received FROM "+
		"(select sender_id as user_id, count(*) as given_kleos, 0 as received_kleos from kleos group by sender_id "+
		"UNION ALL "+
		"select receiver_id as user_id, 0 as given_kleos, count(*) as received_kleos from kleos group by receiver_id) as table1"+
		" group by user_id "+
		"order by sum(given_kleos) desc "+
		"OFFSET ? LIMIT ?", int(pageNumber*pageSize), int(pageSize)).Scan(&adminAllUser)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, int(totalElements), nil
	}
	return adminAllUser, int(totalElements), nil
}

func (r *KleosRepository) GetOneOnOneCount(givenId string, receiverId string, week string) (*KleosGiven, error) {

	var givenStruct KleosGiven
	result := r.gormClient.Table("kleos").
		Select("sender_id, count(*) as count__").
		Where("sender_id = ? AND receiver_id = ? AND week=?", givenId, receiverId, week).
		Group("sender_id").
		Scan(&givenStruct)
	if result.Error != nil {
		return nil, result.Error
	}
	return &givenStruct, nil
}
