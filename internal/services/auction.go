package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"simple-auction-system/internal/models"
	"simple-auction-system/internal/repositories"
)

var (
	ErrInvalidAuction       = errors.New("invalid auction")
	ErrAuctionNotFound      = errors.New("auction not found")
	ErrNotAuctionOwner      = errors.New("only the owner can manage this auction")
	ErrAuctionHasBids       = errors.New("auction cannot be deleted because it has bids")
	ErrAuctionAlreadyClosed = errors.New("auction is already closed")
)

type AuctionService interface {
	CreateAuction(auction *models.Auction) error
	GetAllAuctions() ([]models.Auction, error)
	GetAuctionByID(id uint) (*models.Auction, error)
	DeleteAuction(auctionID uint, userID uint) error
	CloseAuction(auctionID uint, userID uint) (*models.Auction, *models.Bid, error)
}

type auctionService struct {
	repository    repositories.AuctionRepository
	bidRepository repositories.BidRepository
}

func NewAuctionService(
	repository repositories.AuctionRepository,
	bidRepository repositories.BidRepository,
) AuctionService {
	return &auctionService{
		repository:    repository,
		bidRepository: bidRepository,
	}
}

func findAuction(repo repositories.AuctionRepository, id uint) (*models.Auction, error) {
	auction, err := repo.GetAuctionByID(id)
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, ErrAuctionNotFound
	}
	return auction, err
}

func (service *auctionService) CreateAuction(
	auction *models.Auction,
) error {
	auction.Title = strings.TrimSpace(auction.Title)

	if auction.Title == "" {
		return fmt.Errorf(
			"%w: Tiêu đề không được để trống",
			ErrInvalidAuction,
		)
	}

	if auction.StartingPrice <= 0 {
		return fmt.Errorf(
			"%w: giá khởi điểm phải lớn hơn 0",
			ErrInvalidAuction,
		)
	}

	if auction.MinIncrement <= 0 {
		return fmt.Errorf(
			"%w: bước nhảy tối thiểu phải lớn hơn 0",
			ErrInvalidAuction,
		)
	}

	if !auction.EndTime.After(time.Now()) {
		return fmt.Errorf(
			"%w: thời gian kết thúc phải lớn hơn thời gian hiện tại",
			ErrInvalidAuction,
		)
	}

	auction.Status = models.AuctionStatusActive

	return service.repository.CreateAuction(auction)
}

func (service *auctionService) GetAllAuctions() (
	[]models.Auction,
	error,
) {
	return service.repository.GetAllAuctions()
}

func (service *auctionService) GetAuctionByID(
	id uint,
) (*models.Auction, error) {
	return findAuction(service.repository, id)
}

func (service *auctionService) DeleteAuction(
	auctionID uint,
	userID uint,
) error {
	auction, err := findAuction(service.repository, auctionID)
	if err != nil {
		return err
	}

	if auction.OwnerID != userID {
		return ErrNotAuctionOwner
	}

	bidCount, err := service.bidRepository.
		CountBidsByAuctionID(auctionID)

	if err != nil {
		return err
	}

	if bidCount > 0 {
		return ErrAuctionHasBids
	}

	return service.repository.DeleteAuction(auctionID)
}

func (service *auctionService) CloseAuction(
	auctionID uint,
	userID uint,
) (*models.Auction, *models.Bid, error) {
	auction, err := findAuction(service.repository, auctionID)
	if err != nil {
		return nil, nil, err
	}

	if auction.OwnerID != userID {
		return nil, nil, ErrNotAuctionOwner
	}

	if auction.Status == models.AuctionStatusClosed {
		return nil, nil, ErrAuctionAlreadyClosed
	}

	winningBid, err := service.bidRepository.GetHighestBid(auctionID)
	if err != nil && !errors.Is(err, repositories.ErrNotFound) {
		return nil, nil, err
	}

	err = service.repository.UpdateAuctionStatus(
		auctionID,
		models.AuctionStatusClosed,
	)
	if err != nil {
		return nil, nil, err
	}

	auction.Status = models.AuctionStatusClosed

	return auction, winningBid, nil
}
