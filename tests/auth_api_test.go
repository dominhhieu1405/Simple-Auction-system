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

const baseURL = "http://localhost:8080/api"

type apiResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Message  string `json:"message"`
	Token    string `json:"token"`
	Error    string `json:"error"`
}

func sendRequest(
	t *testing.T,
	method string,
	url string,
	token string,
	payload any,
) (int, apiResponse) {
	t.Helper()

	var body io.Reader

	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal request: %v", err)
		}

		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf(
			"send request: %v\nServer đã được chạy chưa?",
			err,
		)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}

	var result apiResponse

	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf(
			"decode response %q: %v",
			string(data),
			err,
		)
	}

	return resp.StatusCode, result
}

func TestAuthAPI(t *testing.T) {
	username := fmt.Sprintf(
		"apitest_%d",
		time.Now().UnixNano(),
	)
	password := "password123"

	var token string

	t.Run("register", func(t *testing.T) {
		status, result := sendRequest(
			t,
			http.MethodPost,
			baseURL+"/auth/register",
			"",
			map[string]string{
				"username": username,
				"password": password,
			},
		)

		if status != http.StatusCreated {
			t.Fatalf(
				"expected 201, got %d; error=%q",
				status,
				result.Error,
			)
		}

		if result.Message != "user registered successfully" {
			t.Errorf(
				"unexpected message: %q",
				result.Message,
			)
		}
	})

	t.Run("login", func(t *testing.T) {
		status, result := sendRequest(
			t,
			http.MethodPost,
			baseURL+"/auth/login",
			"",
			map[string]string{
				"username": username,
				"password": password,
			},
		)

		if status != http.StatusOK {
			t.Fatalf(
				"expected 200, got %d; error=%q",
				status,
				result.Error,
			)
		}

		if result.Token == "" {
			t.Fatal("expected JWT token, got empty token")
		}

		token = result.Token
	})

	t.Run("get current user", func(t *testing.T) {
		if token == "" {
			t.Fatal("login did not return a token")
		}

		status, result := sendRequest(
			t,
			http.MethodGet,
			baseURL+"/users/me",
			token,
			nil,
		)

		if status != http.StatusOK {
			t.Fatalf(
				"expected 200, got %d; error=%q",
				status,
				result.Error,
			)
		}

		if result.ID == 0 {
			t.Error("expected user ID greater than 0")
		}

		if result.Username != username {
			t.Errorf(
				"expected username %q, got %q",
				username,
				result.Username,
			)
		}
	})
}
