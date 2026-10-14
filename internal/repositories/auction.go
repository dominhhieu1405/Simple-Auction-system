package repositories

import (
	"errors"

	"simple-auction-system/internal/models"

	"gorm.io/gorm"
)

type AuctionRepository interface {
	CreateAuction(auction *models.Auction) error
	GetAllAuctions() ([]models.Auction, error)
	GetAuctionByID(id uint) (*models.Auction, error)
	DeleteAuction(id uint) error
	UpdateAuctionStatus(id uint, status string) error
}

type auctionRepository struct {
	database *gorm.DB
}

func NewAuctionRepository(
	database *gorm.DB,
) AuctionRepository {
	return &auctionRepository{
		database: database,
	}
}

func (repository *auctionRepository) CreateAuction(
	auction *models.Auction,
) error {
	return repository.database.Create(auction).Error
}

func (repository *auctionRepository) GetAllAuctions() (
	[]models.Auction,
	error,
) {
	auctions := make([]models.Auction, 0)

	err := repository.database.
		Order("created_at DESC").
		Find(&auctions).
		Error

	return auctions, err
}

func (repository *auctionRepository) GetAuctionByID(
	id uint,
) (*models.Auction, error) {
	var auction models.Auction

	err := repository.database.
		First(&auction, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &auction, nil
}

func (repository *auctionRepository) DeleteAuction(
	id uint,
) error {
	return repository.database.
		Delete(&models.Auction{}, id).
		Error
}

func (repository *auctionRepository) UpdateAuctionStatus(
	id uint,
	status string,
) error {
	return repository.database.
		Model(&models.Auction{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}
