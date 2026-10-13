package repositories

import (
	"simple-auction-system/internal/models"

	"gorm.io/gorm"
)

type AuctionRepository interface {
	CreateAuction(auction *models.Auction) error
	GetAllAuctions() ([]models.Auction, error)
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
