package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
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

// CreateAuction tạo phiên đấu giá.
//
// @Summary Tạo phiên đấu giá
// @Tags Auctions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateAuctionRequest true "Thông tin phiên đấu giá"
// @Success 201 {object} models.Auction
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auctions [post]
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

// GetAllAuctions lấy danh sách phiên đấu giá.
//
// @Summary Danh sách phiên đấu giá
// @Tags Auctions
// @Produce json
// @Success 200 {array} models.Auction
// @Failure 500 {object} ErrorResponse
// @Router /auctions [get]
func (handler *AuctionHandler) GetAllAuctions(
	writer http.ResponseWriter,
	request *http.Request,
) {
	auctions, err := handler.service.GetAllAuctions()
	if err != nil {
		log.Printf("get auctions: %v", err)

		writeError(
			writer,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	writeJSON(
		writer,
		http.StatusOK,
		auctions,
	)
}

// GetAuctionByID lấy chi tiết phiên đấu giá.
//
// @Summary Chi tiết phiên đấu giá
// @Tags Auctions
// @Produce json
// @Param id path int true "Auction ID"
// @Success 200 {object} models.Auction
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auctions/{id} [get]
func (handler *AuctionHandler) GetAuctionByID(
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

	auction, err := handler.service.GetAuctionByID(uint(id))

	switch {
	case err == nil:
		writeJSON(
			writer,
			http.StatusOK,
			auction,
		)

	case errors.Is(err, services.ErrAuctionNotFound):
		writeError(
			writer,
			http.StatusNotFound,
			err.Error(),
		)

	default:
		log.Printf("get auction: %v", err)

		writeError(
			writer,
			http.StatusInternalServerError,
			"internal server error",
		)
	}
}
