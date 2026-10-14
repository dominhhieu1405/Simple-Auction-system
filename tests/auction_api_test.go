package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

const auctionBaseAddress = "http://localhost:8080/api"

type testUserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type testAuctionResponse struct {
	ID            uint      `json:"id"`
	OwnerID       uint      `json:"owner_id"`
	Title         string    `json:"title"`
	StartingPrice int64     `json:"starting_price"`
	MinIncrement  int64     `json:"min_increment"`
	EndTime       time.Time `json:"end_time"`
	Status        string    `json:"status"`
}

type testBidResponse struct {
	ID        uint  `json:"id"`
	AuctionID uint  `json:"auction_id"`
	UserID    uint  `json:"user_id"`
	Amount    int64 `json:"amount"`
}

type testCloseAuctionResponse struct {
	Auction    testAuctionResponse `json:"auction"`
	WinningBid *testBidResponse    `json:"winning_bid"`
}

func sendIntegrationRequest(
	test *testing.T,
	method string,
	address string,
	token string,
	requestData any,
) (int, []byte) {
	test.Helper()

	var requestBody io.Reader

	if requestData != nil {
		data, err := json.Marshal(requestData)
		if err != nil {
			test.Fatal(err)
		}

		requestBody = bytes.NewReader(data)
	}

	request, err := http.NewRequest(
		method,
		address,
		requestBody,
	)
	if err != nil {
		test.Fatal(err)
	}

	if requestData != nil {
		request.Header.Set(
			"Content-Type",
			"application/json",
		)
	}

	if token != "" {
		request.Header.Set(
			"Authorization",
			"Bearer "+token,
		)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		test.Fatalf(
			"Không gửi được request: %v\nServer đã chạy chưa?",
			err,
		)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		test.Fatal(err)
	}

	return response.StatusCode, responseBody
}

func createAuthenticatedTestUser(
	test *testing.T,
	namePrefix string,
) (uint, string) {
	test.Helper()

	username := fmt.Sprintf(
		"%s_%d",
		namePrefix,
		time.Now().UnixNano(),
	)
	password := "password123"

	status, responseBody := sendIntegrationRequest(
		test,
		http.MethodPost,
		auctionBaseAddress+"/auth/register",
		"",
		map[string]string{
			"username": username,
			"password": password,
		},
	)

	if status != http.StatusCreated {
		test.Fatalf(
			"Đăng ký thất bại, status=%d, body=%s",
			status,
			responseBody,
		)
	}

	status, responseBody = sendIntegrationRequest(
		test,
		http.MethodPost,
		auctionBaseAddress+"/auth/login",
		"",
		map[string]string{
			"username": username,
			"password": password,
		},
	)

	if status != http.StatusOK {
		test.Fatalf(
			"Đăng nhập thất bại, status=%d, body=%s",
			status,
			responseBody,
		)
	}

	var loginResult testUserResponse

	if err := json.Unmarshal(
		responseBody,
		&loginResult,
	); err != nil {
		test.Fatal(err)
	}

	if loginResult.Token == "" {
		test.Fatal("Đăng nhập không trả về token")
	}

	status, responseBody = sendIntegrationRequest(
		test,
		http.MethodGet,
		auctionBaseAddress+"/users/me",
		loginResult.Token,
		nil,
	)

	if status != http.StatusOK {
		test.Fatalf(
			"Lấy user thất bại, status=%d, body=%s",
			status,
			responseBody,
		)
	}

	var userResult testUserResponse

	if err := json.Unmarshal(
		responseBody,
		&userResult,
	); err != nil {
		test.Fatal(err)
	}

	return userResult.ID, loginResult.Token
}

func TestAuctionFlow(test *testing.T) {
	ownerID, ownerToken := createAuthenticatedTestUser(
		test,
		"owner",
	)

	bidderID, bidderToken := createAuthenticatedTestUser(
		test,
		"bidder",
	)

	var auction testAuctionResponse

	test.Run("create auction", func(test *testing.T) {
		status, responseBody := sendIntegrationRequest(
			test,
			http.MethodPost,
			auctionBaseAddress+"/auctions",
			ownerToken,
			map[string]any{
				"title":          "Test auction",
				"description":    "Auction created by API test",
				"starting_price": int64(1000000),
				"min_increment":  int64(50000),
				"end_time":       time.Now().Add(2 * time.Hour),
			},
		)

		if status != http.StatusCreated {
			test.Fatalf(
				"expected 201, got %d; body=%s",
				status,
				responseBody,
			)
		}

		if err := json.Unmarshal(
			responseBody,
			&auction,
		); err != nil {
			test.Fatal(err)
		}

		if auction.ID == 0 {
			test.Fatal("Auction ID bằng 0")
		}

		if auction.OwnerID != ownerID {
			test.Fatalf(
				"expected owner %d, got %d",
				ownerID,
				auction.OwnerID,
			)
		}

		if auction.Status != "active" {
			test.Fatalf(
				"expected active, got %q",
				auction.Status,
			)
		}
	})

	test.Run("get auctions", func(test *testing.T) {
		status, responseBody := sendIntegrationRequest(
			test,
			http.MethodGet,
			auctionBaseAddress+"/auctions",
			"",
			nil,
		)

		if status != http.StatusOK {
			test.Fatalf(
				"expected 200, got %d; body=%s",
				status,
				responseBody,
			)
		}

		var auctions []testAuctionResponse

		if err := json.Unmarshal(
			responseBody,
			&auctions,
		); err != nil {
			test.Fatal(err)
		}

		if len(auctions) == 0 {
			test.Fatal("Danh sách auction đang trống")
		}
	})

	test.Run("get auction detail", func(test *testing.T) {
		address := fmt.Sprintf(
			"%s/auctions/%d",
			auctionBaseAddress,
			auction.ID,
		)

		status, responseBody := sendIntegrationRequest(
			test,
			http.MethodGet,
			address,
			"",
			nil,
		)

		if status != http.StatusOK {
			test.Fatalf(
				"expected 200, got %d; body=%s",
				status,
				responseBody,
			)
		}

		var result testAuctionResponse

		if err := json.Unmarshal(
			responseBody,
			&result,
		); err != nil {
			test.Fatal(err)
		}

		if result.ID != auction.ID {
			test.Fatalf(
				"expected auction %d, got %d",
				auction.ID,
				result.ID,
			)
		}
	})

	var bid testBidResponse

	test.Run("place bid", func(test *testing.T) {
		address := fmt.Sprintf(
			"%s/auctions/%d/bids",
			auctionBaseAddress,
			auction.ID,
		)

		status, responseBody := sendIntegrationRequest(
			test,
			http.MethodPost,
			address,
			bidderToken,
			map[string]int64{
				"amount": 1000000,
			},
		)

		if status != http.StatusCreated {
			test.Fatalf(
				"expected 201, got %d; body=%s",
				status,
				responseBody,
			)
		}

		if err := json.Unmarshal(
			responseBody,
			&bid,
		); err != nil {
			test.Fatal(err)
		}

		if bid.UserID != bidderID {
			test.Fatalf(
				"expected bidder %d, got %d",
				bidderID,
				bid.UserID,
			)
		}
	})

	test.Run("get auction bids", func(test *testing.T) {
		address := fmt.Sprintf(
			"%s/auctions/%d/bids",
			auctionBaseAddress,
			auction.ID,
		)

		status, responseBody := sendIntegrationRequest(
			test,
			http.MethodGet,
			address,
			"",
			nil,
		)

		if status != http.StatusOK {
			test.Fatalf(
				"expected 200, got %d; body=%s",
				status,
				responseBody,
			)
		}

		var bids []testBidResponse

		if err := json.Unmarshal(
			responseBody,
			&bids,
		); err != nil {
			test.Fatal(err)
		}

		if len(bids) != 1 {
			test.Fatalf(
				"expected 1 bid, got %d",
				len(bids),
			)
		}
	})

	test.Run("get current user bids", func(test *testing.T) {
		status, responseBody := sendIntegrationRequest(
			test,
			http.MethodGet,
			auctionBaseAddress+"/users/me/bids",
			bidderToken,
			nil,
		)

		if status != http.StatusOK {
			test.Fatalf(
				"expected 200, got %d; body=%s",
				status,
				responseBody,
			)
		}

		var bids []testBidResponse

		if err := json.Unmarshal(
			responseBody,
			&bids,
		); err != nil {
			test.Fatal(err)
		}

		if len(bids) == 0 {
			test.Fatal("Không tìm thấy lịch sử bid")
		}
	})

	test.Run("close auction", func(test *testing.T) {
		address := fmt.Sprintf(
			"%s/auctions/%d/close",
			auctionBaseAddress,
			auction.ID,
		)

		status, responseBody := sendIntegrationRequest(
			test,
			http.MethodPost,
			address,
			ownerToken,
			nil,
		)

		if status != http.StatusOK {
			test.Fatalf(
				"expected 200, got %d; body=%s",
				status,
				responseBody,
			)
		}

		var result testCloseAuctionResponse

		if err := json.Unmarshal(
			responseBody,
			&result,
		); err != nil {
			test.Fatal(err)
		}

		if result.Auction.Status != "closed" {
			test.Fatalf(
				"expected closed, got %q",
				result.Auction.Status,
			)
		}

		if result.WinningBid == nil {
			test.Fatal("Không có winning bid")
		}

		if result.WinningBid.UserID != bidderID {
			test.Fatalf(
				"expected winner %d, got %d",
				bidderID,
				result.WinningBid.UserID,
			)
		}
	})

	test.Run("reject bid after close", func(test *testing.T) {
		address := fmt.Sprintf(
			"%s/auctions/%d/bids",
			auctionBaseAddress,
			auction.ID,
		)

		status, responseBody := sendIntegrationRequest(
			test,
			http.MethodPost,
			address,
			bidderToken,
			map[string]int64{
				"amount": 1100000,
			},
		)

		if status != http.StatusConflict {
			test.Fatalf(
				"expected 409, got %d; body=%s",
				status,
				responseBody,
			)
		}
	})

	test.Run("reject delete auction with bids", func(test *testing.T) {
		address := fmt.Sprintf(
			"%s/auctions/%d",
			auctionBaseAddress,
			auction.ID,
		)

		status, responseBody := sendIntegrationRequest(
			test,
			http.MethodDelete,
			address,
			ownerToken,
			nil,
		)

		if status != http.StatusConflict {
			test.Fatalf(
				"expected 409, got %d; body=%s",
				status,
				responseBody,
			)
		}
	})
}
