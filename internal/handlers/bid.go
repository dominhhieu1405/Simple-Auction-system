package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"simple-auction-system/internal/middleware"
	"simple-auction-system/internal/services"
)

type BidHandler struct {
	service services.BidService
}

func NewBidHandler(
	service services.BidService,
) *BidHandler {
	return &BidHandler{
		service: service,
	}
}

type PlaceBidRequest struct {
	Amount int64 `json:"amount"`
}

// PlaceBid đặt giá cho phiên đấu giá.
//
// @Summary Đặt giá
// @Tags Bids
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Auction ID"
// @Param request body PlaceBidRequest true "Số tiền đặt giá"
// @Success 201 {object} models.Bid
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auctions/{id}/bids [post]
func (handler *BidHandler) PlaceBid(
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

	idText := request.PathValue("id")

	id, err := strconv.ParseUint(idText, 10, 64)
	if err != nil || id == 0 {
		writeError(
			writer,
			http.StatusBadRequest,
			"invalid auction id",
		)
		return
	}

	var requestData PlaceBidRequest

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

	bid, err := handler.service.PlaceBid(
		uint(id),
		userID,
		requestData.Amount,
	)

	switch {
	case err == nil:
		writeJSON(
			writer,
			http.StatusCreated,
			bid,
		)

	case errors.Is(err, services.ErrInvalidBid),
		errors.Is(err, services.ErrBidTooLow):
		writeError(
			writer,
			http.StatusBadRequest,
			err.Error(),
		)

	case errors.Is(err, services.ErrAuctionNotFound):
		writeError(
			writer,
			http.StatusNotFound,
			err.Error(),
		)

	case errors.Is(err, services.ErrAuctionNotActive),
		errors.Is(err, services.ErrAuctionEnded):
		writeError(
			writer,
			http.StatusConflict,
			err.Error(),
		)

	default:
		log.Printf("place bid: %v", err)

		writeError(
			writer,
			http.StatusInternalServerError,
			"internal server error",
		)
	}

	
}
// GetBidsByAuctionID lấy lịch sử đặt giá.
//
// @Summary Lịch sử đặt giá
// @Tags Bids
// @Produce json
// @Param id path int true "Auction ID"
// @Success 200 {array} models.Bid
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auctions/{id}/bids [get]
func (handler *BidHandler) GetBidsByAuctionID(
	writer http.ResponseWriter,
	request *http.Request,
) {
	idText := request.PathValue("id")

	id, err := strconv.ParseUint(idText, 10, 64)
	if err != nil || id == 0 {
		writeError(
			writer,
			http.StatusBadRequest,
			"invalid auction id",
		)
		return
	}

	bids, err := handler.service.
		GetBidsByAuctionID(uint(id))

	switch {
	case err == nil:
		writeJSON(
			writer,
			http.StatusOK,
			bids,
		)

	case errors.Is(err, services.ErrAuctionNotFound):
		writeError(
			writer,
			http.StatusNotFound,
			err.Error(),
		)

	default:
		log.Printf("get auction bids: %v", err)

		writeError(
			writer,
			http.StatusInternalServerError,
			"internal server error",
		)
	}
}