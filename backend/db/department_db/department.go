package departmentDb

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DepartmentRepository struct {
	logger     *zap.Logger
	gormClient *gorm.DB
}

func NewDepartmentRepository(logger *zap.Logger, gormClient *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{
		logger:     logger,
		gormClient: gormClient,
	}
}

func (h *DepartmentRepository) GetDepartmentData(id int) (*DepartmentEntity, error) {

	if h == nil || h.gormClient == nil {
		return nil, errors.New("gormClient is nil")
	}
	var departmentEntity DepartmentEntity
	result := h.gormClient.Table("departments").Where("id = ?", id).Scan(&departmentEntity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &departmentEntity, nil
}

func (h *DepartmentRepository) GetDepartmentDataFromName(departmentName string) (*DepartmentEntity, error) {

	if h == nil || h.gormClient == nil {
		return nil, errors.New("gormClient is nil")
	}

	var departmentEntity DepartmentEntity

	result := h.gormClient.Table("departments").Where("name = ?", departmentName).Scan(&departmentEntity)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &departmentEntity, nil
}

func (h *DepartmentRepository) CreateDepartment(departmentName string) (*DepartmentEntity, error) {
	department := DepartmentEntity{Name: departmentName}
	err := h.gormClient.Create(&department).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

func (h *DepartmentRepository) GetDepartmentDataById(id int) (*DepartmentEntity, error) {

	if h == nil || h.gormClient == nil {
		return nil, errors.New("gormClient is nil")
	}
	var departmentEntity DepartmentEntity
	result := h.gormClient.Table("departments").Where("id = ?", id).Scan(&departmentEntity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &departmentEntity, nil
}
