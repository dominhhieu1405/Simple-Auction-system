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

type authResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	Error   string `json:"error"`
}

func sendJSON(
	t *testing.T,
	method string,
	url string,
	payload any,
) (*http.Response, authResponse) {
	t.Helper()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("send request: %v\nServer đã được chạy chưa?", err)
	}
	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	var result authResponse
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("decode response %q: %v", string(data), err)
	}

	return resp, result
}

func TestRegisterAndLogin(t *testing.T) {
	// Tạo username khác nhau sau mỗi lần chạy để tránh bị trùng.
	username := fmt.Sprintf("apitest_%d", time.Now().UnixNano())
	password := "password123"

	t.Run("register successfully", func(t *testing.T) {
		resp, result := sendJSON(
			t,
			http.MethodPost,
			baseURL+"/auth/register",
			map[string]string{
				"username": username,
				"password": password,
			},
		)

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf(
				"expected status 201, got %d; error=%q",
				resp.StatusCode,
				result.Error,
			)
		}

		if result.Message != "user registered successfully" {
			t.Errorf(
				"unexpected message: got %q",
				result.Message,
			)
		}
	})

	t.Run("login successfully", func(t *testing.T) {
		resp, result := sendJSON(
			t,
			http.MethodPost,
			baseURL+"/auth/login",
			map[string]string{
				"username": username,
				"password": password,
			},
		)

		if resp.StatusCode != http.StatusOK {
			t.Fatalf(
				"expected status 200, got %d; error=%q",
				resp.StatusCode,
				result.Error,
			)
		}

		if result.Token == "" {
			t.Fatal("expected JWT token, got empty token")
		}
	})
}
