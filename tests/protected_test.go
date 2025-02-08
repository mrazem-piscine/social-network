package tests

import (
	"net/http"
	"net/http/httptest"
	"social-network/internal/middlewares"
	"testing"
)

func TestProtectedRouteWithoutAuth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
	rec := httptest.NewRecorder()

	protectedHandler := middlewares.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	protectedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rec.Code)
	}
}
