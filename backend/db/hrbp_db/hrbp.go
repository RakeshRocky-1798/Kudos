package hrbpDb

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HrbpRepository struct {
	logger     *zap.Logger
	gormClient *gorm.DB
}

func NewHrbpRepository(logger *zap.Logger, gormClient *gorm.DB) *HrbpRepository {
	return &HrbpRepository{
		logger:     logger,
		gormClient: gormClient,
	}
}

func (h *HrbpRepository) GetHrbpData(id int) (*HrbpEntity, error) {

	if h == nil || h.gormClient == nil {
		return nil, errors.New("gormClient is nil")
	}
	var hrbpEntity HrbpEntity
	result := h.gormClient.Table("hrbp").Where("id = ?", id).Scan(&hrbpEntity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &hrbpEntity, nil
}

func (h *HrbpRepository) GetHrbpDataFromName(hrbpName string) (*HrbpEntity, error) {

	if h == nil || h.gormClient == nil {
		return nil, errors.New("gormClient is nil")
	}

	var hrbpEntity HrbpEntity

	result := h.gormClient.Table("hrbp").Where("name = ?", hrbpName).Scan(&hrbpEntity)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &hrbpEntity, nil
}

func (h *HrbpRepository) CreateHrbp(hrbpName string, email string) (*HrbpEntity, error) {

	hrbp := HrbpEntity{Name: hrbpName, Email: email}
	err := h.gormClient.Create(&hrbp).Error
	if err != nil {
		return nil, err
	}
	return &hrbp, nil
}

func (h *HrbpRepository) GetHrbpDataById(id int) (*HrbpEntity, error) {

	if h == nil || h.gormClient == nil {
		return nil, errors.New("gormClient is nil")
	}
	var hrbpEntity HrbpEntity

	result := h.gormClient.Table("hrbp").Where("id = ?", id).Scan(&hrbpEntity)

	if result.Error != nil {
		return nil, result.Error
	}

	return &hrbpEntity, nil
}
