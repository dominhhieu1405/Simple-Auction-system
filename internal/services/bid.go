package services

import (
	"errors"
	"fmt"
	"time"

	"simple-auction-system/internal/models"
	"simple-auction-system/internal/repositories"
)

var (
	ErrInvalidBid       = errors.New("invalid bid")
	ErrBidTooLow        = errors.New("bid amount is too low")
	ErrAuctionNotActive = errors.New("auction is not active")
	ErrAuctionEnded     = errors.New("auction has ended")
)

type BidService interface {
	PlaceBid(
		auctionID uint,
		userID uint,
		amount int64,
	) (*models.Bid, error)

	GetBidsByAuctionID(
		auctionID uint,
	) ([]models.Bid, error)

	GetBidsByUserID(
		userID uint,
	) ([]models.Bid, error)
}

type bidService struct {
	auctionRepository repositories.AuctionRepository
	bidRepository     repositories.BidRepository
}

func NewBidService(
	auctionRepository repositories.AuctionRepository,
	bidRepository repositories.BidRepository,
) BidService {
	return &bidService{
		auctionRepository: auctionRepository,
		bidRepository:     bidRepository,
	}
}

func (service *bidService) GetBidsByAuctionID(
	auctionID uint,
) ([]models.Bid, error) {
	if _, err := findAuction(service.auctionRepository, auctionID); err != nil {
		return nil, err
	}

	return service.bidRepository.
		GetBidsByAuctionID(auctionID)
}

func (service *bidService) PlaceBid(
	auctionID uint,
	userID uint,
	amount int64,
) (*models.Bid, error) {
	if amount <= 0 {
		return nil, ErrInvalidBid
	}

	auction, err := findAuction(service.auctionRepository, auctionID)
	if err != nil {
		return nil, err
	}

	if auction.Status != models.AuctionStatusActive {
		return nil, ErrAuctionNotActive
	}

	if !auction.EndTime.After(time.Now()) {
		return nil, ErrAuctionEnded
	}

	highestBid, err := service.bidRepository.
		GetHighestBid(auctionID)

	minimumAmount := auction.StartingPrice

	if err == nil {
		minimumAmount = highestBid.Amount +
			auction.MinIncrement
	} else if !errors.Is(err, repositories.ErrNotFound) {
		return nil, err
	}

	if amount < minimumAmount {
		return nil, fmt.Errorf(
			"%w: minimum amount is %d",
			ErrBidTooLow,
			minimumAmount,
		)
	}

	bid := &models.Bid{
		AuctionID: auctionID,
		UserID:    userID,
		Amount:    amount,
	}

	if err := service.bidRepository.CreateBid(bid); err != nil {
		return nil, err
	}

	return bid, nil
}

func (service *bidService) GetBidsByUserID(
	userID uint,
) ([]models.Bid, error) {
	return service.bidRepository.
		GetBidsByUserID(userID)
}
