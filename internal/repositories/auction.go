package repositories

import (
	"simple-auction-system/internal/models"

	"gorm.io/gorm"
)

type AuctionRepository interface {
	CreateAuction(auction *models.Auction) error
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
