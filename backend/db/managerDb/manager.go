package managersDb

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ManagerRepository struct {
	logger     *zap.Logger
	gormClient *gorm.DB
}

func NewManagerRepository(logger *zap.Logger, gormClient *gorm.DB) *ManagerRepository {
	return &ManagerRepository{
		logger:     logger,
		gormClient: gormClient,
	}
}

func (m *ManagerRepository) GetManagerDataFromId(id int) (*ManagerEntity, error) {
	var managerEntity ManagerEntity
	result := m.gormClient.Table("managers").Select("*").
		Where("id=?", id).Scan(&managerEntity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &managerEntity, nil
}

func (m *ManagerRepository) GetManagerDataFromEmail(email string) (*ManagerEntity, error) {
	var managerEntity ManagerEntity

	result := m.gormClient.Table("managers").Select("*").
		Where("email=?", email).Scan(&managerEntity)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &managerEntity, nil
}

func (m *ManagerRepository) CreateManager(managerEntity *ManagerEntity) (*ManagerEntity, error) {

	result := m.gormClient.Table("managers").Create(managerEntity)
	if result.Error != nil {
		return nil, result.Error
	}

	return managerEntity, nil
}
