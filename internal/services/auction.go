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
	ErrInvalidAuction  = errors.New("invalid auction")
	ErrAuctionNotFound = errors.New("auction not found")
)

type AuctionService interface {
	CreateAuction(auction *models.Auction) error
	GetAllAuctions() ([]models.Auction, error)
	GetAuctionByID(id uint) (*models.Auction, error)
}

type auctionService struct {
	repository repositories.AuctionRepository
}

func NewAuctionService(
	repository repositories.AuctionRepository,
) AuctionService {
	return &auctionService{
		repository: repository,
	}
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
	auction, err := service.repository.GetAuctionByID(id)

	if errors.Is(err, repositories.ErrNotFound) {
		return nil, ErrAuctionNotFound
	}

	if err != nil {
		return nil, err
	}

	return auction, nil
}
