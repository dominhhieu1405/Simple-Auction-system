package repositories

import (
	"errors"

	"simple-auction-system/internal/models"

	"gorm.io/gorm"
)

type BidRepository interface {
	CreateBid(bid *models.Bid) error
	GetHighestBid(auctionID uint) (*models.Bid, error)
	GetBidsByAuctionID(auctionID uint) ([]models.Bid, error)
}
func (repository *bidRepository) GetBidsByAuctionID(
	auctionID uint,
) ([]models.Bid, error) {
	bids := make([]models.Bid, 0)

	err := repository.database.
		Where("auction_id = ?", auctionID).
		Order("created_at DESC").
		Find(&bids).
		Error

	return bids, err
}
type bidRepository struct {
	database *gorm.DB
}

func NewBidRepository(
	database *gorm.DB,
) BidRepository {
	return &bidRepository{
		database: database,
	}
}

func (repository *bidRepository) CreateBid(
	bid *models.Bid,
) error {
	return repository.database.Create(bid).Error
}

func (repository *bidRepository) GetHighestBid(
	auctionID uint,
) (*models.Bid, error) {
	var bid models.Bid

	err := repository.database.
		Where("auction_id = ?", auctionID).
		Order("amount DESC").
		First(&bid).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &bid, nil
}