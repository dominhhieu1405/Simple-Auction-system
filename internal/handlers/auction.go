package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"simple-auction-system/internal/middleware"
	"simple-auction-system/internal/models"
	"simple-auction-system/internal/services"
)

type AuctionHandler struct {
	service services.AuctionService
}

func NewAuctionHandler(
	service services.AuctionService,
) *AuctionHandler {
	return &AuctionHandler{
		service: service,
	}
}

type CreateAuctionRequest struct {
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	StartingPrice int64     `json:"starting_price"`
	MinIncrement  int64     `json:"min_increment"`
	EndTime       time.Time `json:"end_time"`
}

func (handler *AuctionHandler) CreateAuction(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userID, ok := middleware.GetUserID(request)
	if !ok {
		writeError(
			writer,
			http.StatusUnauthorized,
			"unauthorized",
		)
		return
	}

	var requestData CreateAuctionRequest

	if err := decodeJSON(
		writer,
		request,
		&requestData,
	); err != nil {
		writeError(
			writer,
			http.StatusBadRequest,
			"invalid request payload",
		)
		return
	}

	auction := models.Auction{
		OwnerID:       userID,
		Title:         requestData.Title,
		Description:   requestData.Description,
		StartingPrice: requestData.StartingPrice,
		MinIncrement:  requestData.MinIncrement,
		EndTime:       requestData.EndTime,
	}

	err := handler.service.CreateAuction(&auction)

	switch {
	case err == nil:
		writeJSON(
			writer,
			http.StatusCreated,
			auction,
		)

	case errors.Is(err, services.ErrInvalidAuction):
		writeError(
			writer,
			http.StatusBadRequest,
			err.Error(),
		)

	default:
		log.Printf("create auction: %v", err)

		writeError(
			writer,
			http.StatusInternalServerError,
			"internal server error",
		)
	}
}
