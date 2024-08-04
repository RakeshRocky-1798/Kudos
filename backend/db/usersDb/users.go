package usersDb

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type UserRepository struct {
	logger     *zap.Logger
	gormClient *gorm.DB
}

func NewUserRepository(logger *zap.Logger, gormClient *gorm.DB) *UserRepository {
	return &UserRepository{
		logger:     logger,
		gormClient: gormClient,
	}
}

func (u *UserRepository) AddUser(request *UserRequest) (*UserEntity, error) {

	userEntity := &UserEntity{
		SlackUserId:   request.SlackUserId,
		Name:          request.UserName,
		Email:         request.Email,
		SlackImageUrl: request.SlackImageUrl,
		RealName:      request.RealName,
		GivenCount:    request.GivenCount,
		ReceivedCount: request.ReceivedCount,
		ManagerId:     0,
		HrbpId:        0,
		DepartmentId:  0,
	}

	result := u.gormClient.Create(userEntity)
	if result.Error != nil {
		return nil, result.Error
	}

	return userEntity, nil
}

func (u *UserRepository) UpdateUser(updatedUser *UserEntity) (*UserEntity, error) {

	result := u.gormClient.Save(updatedUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return updatedUser, nil
}

func (u *UserRepository) GetAllUserData(email string) (*[]UserEmailIdCombo, error) {

	var allUserData []UserEmailIdCombo
	result := u.gormClient.Table("users").Select("id, email").
		Order("id")
	if email != "" {
		result = result.Where("email != ?", email)
	}
	result = result.Scan(&allUserData)
	if result.Error != nil {
		return nil, result.Error
	}
	return &allUserData, nil
}

func (u *UserRepository) GetUserDataFromGmail(gmail string) (*UserEntity, error) {

	if u == nil || u.gormClient == nil {
		u.logger.Error("[GetUserDataFromGmail]Repository or gormClient is nil")
		return nil, errors.New("[GetUserDataFromGmail]repository or gormClient is nil")
	}

	userName := strings.Split(gmail, "@")[0]
	var userEntity UserEntity
	result := u.gormClient.Table("users").Select("*").
		Where("user_name=?", userName).Scan(&userEntity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userEntity, nil
}

func (u *UserRepository) GetSenderIdFromEmail(email string) (int, error) {

	if u == nil || u.gormClient == nil {
		u.logger.Error("[GetUserDataFromGmail]Repository or gormClient is nil")
		return -1, errors.New("[GetUserDataFromGmail]repository or gormClient is nil")
	}

	var userId int
	var userEntity UserEntity

	result := u.gormClient.Table("users").Select("*").
		Where("email=?", email).Scan(&userEntity)

	if result.Error != nil {
		return -1, result.Error
	}

	userId = int(userEntity.ID)

	return userId, nil
}

func (u *UserRepository) GetUserInfoFromId(userId string) (*CustomUserInfo, error) {
	if u == nil || u.gormClient == nil {
		u.logger.Error("[GetUserInfoFromId]Repository or gormClient is nil")
		return nil, errors.New("[GetUserInfoFromId]repository or gormClient is nil")
	}

	newUser, _ := strconv.Atoi(userId)

	var userInfo CustomUserInfo
	result := u.gormClient.Table("users").
		Select("id, email, real_name, slack_image_url").
		Where("id=?", newUser).Scan(&userInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userInfo, nil
}

func (u *UserRepository) GetUserIdFromSlackId(slackId string) (*UserIdInfo, error) {
	if u == nil || u.gormClient == nil {
		u.logger.Error("[GetUserIdFromSlackId]Repository or gormClient is nil")
		return nil, errors.New("[GetUserIdFromSlackId]repository or gormClient is nil")
	}

	var userInfo UserIdInfo
	result := u.gormClient.Table("users").
		Select("*").
		Where("slack_user_id=?", slackId).Scan(&userInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userInfo, nil
}

func (u *UserRepository) GetSlackIdFromId(id int) (*UserIdInfo, error) {
	if u == nil || u.gormClient == nil {
		u.logger.Error("[GetUserNameFromId]Repository or gormClient is nil")
		return nil, errors.New("[GetUserNameFromId]repository or gormClient is nil")
	}
	var userInfo UserIdInfo
	result := u.gormClient.Table("users").
		Select("id, slack_user_id, email, manager_id, slack_image_url, real_name").
		Where("id=?", id).Scan(&userInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userInfo, nil
}

func (u *UserRepository) SetSlackIdFromWeb(email string, slackId string) error {
	if u == nil || u.gormClient == nil {
		u.logger.Error("[GetUserNameFromId]Repository or gormClient is nil")
		return errors.New("[GetUserNameFromId]repository or gormClient is nil")
	}
	result := u.gormClient.Table("users").
		Where("email=?", email).
		Update("slack_user_id", slackId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserRepository) GetUserCount(userId string) (*UserCount, error) {

	if u == nil || u.gormClient == nil {
		u.logger.Error("[UpdateUserReceivedCount]Repository or gormClient is nil")
		return nil, errors.New("[UpdateUserReceivedCount]repository or gormClient is nil")
	}
	var userCount UserCount
	result := u.gormClient.Table("users").
		Select("id, slack_user_id, given_count, received_count").
		Where("id=?", userId).Scan(&userCount)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userCount, nil
}

func (u *UserRepository) setMissingInfo(email string) {
	if u == nil || u.gormClient == nil {
		u.logger.Error("[GetUserNameFromId]Repository or gormClient is nil")
		return
	}
	result := u.gormClient.Table("users").
		Where("email=?", email).
		Update("slack_user_id", "nil")
	if result.Error != nil {
		u.logger.Error("Error while setting the missing info", zap.Error(result.Error))
		return
	}
}

func (u *UserRepository) GetUsersWithSlackId() (*[]UserSlackCombo, error) {

	if u == nil || u.gormClient == nil {
		u.logger.Error("[GetUserNameFromId]Repository or gormClient is nil")
		return nil, errors.New("[GetUserNameFromId]repository or gormClient is nil")
	}
	var userSlackInfo []UserSlackCombo
	result := u.gormClient.Table("users").
		Select("id, slack_user_id").
		Where("slack_user_id !=?", "nil").Scan(&userSlackInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userSlackInfo, nil
}

func (u *UserRepository) CreateUser(userEntity *UserEntity) (*UserEntity, error) {

	result := u.gormClient.Create(userEntity)
	if result.Error != nil {
		return nil, result.Error
	}

	return userEntity, nil
}

func (u *UserRepository) GetUserGivenAndReceivedCount(userId int) (*UserGivenReceivedCount, error) {

	if u == nil || u.gormClient == nil {
		u.logger.Error("[GetUserNameFromId]Repository or gormClient is nil")
		return nil, errors.New("[GetUserNameFromId]repository or gormClient is nil")
	}

	var userGivenReceivedCount UserGivenReceivedCount
	result := u.gormClient.Table("users").
		Select("id, given_count, received_count").
		Where("id=?", userId).Scan(&userGivenReceivedCount)

	if result.Error != nil {
		return nil, result.Error
	}
	return &userGivenReceivedCount, nil
}

func (u *UserRepository) UpdateUserGivenCount(userId string) error {

	if u == nil || u.gormClient == nil {
		u.logger.Error("[UpdateUserReceivedCount]Repository or gormClient is nil")
		return errors.New("[UpdateUserReceivedCount]repository or gormClient is nil")
	}
	sql := `UPDATE users SET given_count = given_count + 1, updated_at= now() WHERE id = ?`
	result := u.gormClient.Exec(sql, userId)
	if result.Error != nil {
		u.logger.Error("[UpdateUserReceivedCount]Error while adding user value")
		return errors.New("[UpdateUserReceivedCount]Error while adding user value")
	}
	return nil
}

func (u *UserRepository) UpdateUserReceivedCount(userId string) error {

	if u == nil || u.gormClient == nil {
		u.logger.Error("[UpdateUserReceivedCount]Repository or gormClient is nil")
		return errors.New("[UpdateUserReceivedCount]repository or gormClient is nil")
	}
	sql := `UPDATE users SET received_count = received_count + 1, updated_at= now() WHERE id = ?`
	result := u.gormClient.Exec(sql, userId)
	if result.Error != nil {
		u.logger.Error("[UpdateUserReceivedCount]Error while adding user value")
		return errors.New("[UpdateUserReceivedCount]Error while adding user value")
	}
	return nil
}

func (u *UserRepository) GetAllUserCount() (*[]UserGivenReceivedCount, error) {

	if u == nil || u.gormClient == nil {
		u.logger.Error("[UpdateUserReceivedCount]Repository or gormClient is nil")
		return nil, errors.New("[GetUserNameFromId]repository or gormClient is nil")
	}
	var userData []UserGivenReceivedCount
	result := u.gormClient.Table("users").
		Select("id, given_count, received_count").
		Order("given_count DESC").
		Scan(&userData)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userData, nil
}

func (u *UserRepository) CheckIfUserExists(id int) (bool, error) {
	if u == nil || u.gormClient == nil {
		u.logger.Error("[UpdateUserReceivedCount]Repository or gormClient is nil")
		return false, errors.New("[GetUserNameFromId]repository or gormClient is nil")
	}

	var count int64

	result := u.gormClient.Table("users").Where("id=?", id).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (u *UserRepository) GetAdminAdvData(id string) (*AdvUserEntity, error) {
	if u == nil || u.gormClient == nil {
		u.logger.Error("[GetAdminAdvData]Repository or gormClient is nil")
		return nil, errors.New("[GetAdminAdvData]repository or gormClient is nil")
	}
	newUser, _ := strconv.Atoi(id)

	var adminEntity AdvUserEntity
	result := u.gormClient.Table("users").Select("*").
		Where("id=?", newUser).Scan(&adminEntity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &adminEntity, nil
}

func (u *UserRepository) IsUserExistAlready(email string) (bool, error) {
	if u == nil || u.gormClient == nil {
		u.logger.Error("[isUserExistAlready]Repository or gormClient is nil")
		return false, errors.New("[isUserExistAlready]repository or gormClient is nil")
	}

	var count int64
	result := u.gormClient.Table("users").Where("email=?", email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (u *UserRepository) GetUserDataIfPresent(email string) (*UserEntity, error) {

	var userEntity UserEntity

	result := u.gormClient.Table("users").Select("*").Joins("JOIN managers ON users.email = ?", email).
		Where("email != ?", email).Scan(&userEntity)

	if result.Error != nil {
		return nil, result.Error
	}
	return &userEntity, nil
}

func (u *UserRepository) GetUserDataWithForeignJoins(email string) (*UserWithForeignData, error) {

	var userWithForeignData UserWithForeignData

	result := u.gormClient.Table("users").
		Select("*").
		Joins("LEFT JOIN managers ON users.manager_id = managers.id").
		Joins("LEFT JOIN departments ON users.department_id = departments.id").
		Joins("LEFT JOIN hrbp ON users.hrbp_id = hrbp.id").
		Where("users.email = ?", email).
		Scan(&userWithForeignData)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &userWithForeignData, nil
}
