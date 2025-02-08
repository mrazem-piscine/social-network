package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"social-network/internal/handlers"
	"testing"
)

func TestUserRegistration(t *testing.T) {
	reqBody, _ := json.Marshal(map[string]string{
		"nickname": "testuser",
		"email":    "test@example.com",
		"password": "testpass",
	})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handlers.RegisterUser(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rec.Code)
	}
}

func TestUserLogin(t *testing.T) {
	reqBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "testpass",
	})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handlers.LoginUser(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}
