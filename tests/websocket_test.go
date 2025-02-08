package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	websockets "social-network/internal/websocket"
	"github.com/gorilla/websocket"
)

func TestWebSocketConnection(t *testing.T) {
	// ✅ Define WebSocket upgrader inside test function
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

	// Create a test WebSocket server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade WebSocket connection: %v", err)
		}
		manager := websockets.NewWebSocketManager()
		go manager.HandleConnection(conn, "testuser", nil)
	}))
	defer server.Close()

	// ✅ Use `websocket.DefaultDialer`
	url := "ws" + server.URL[4:] + "/ws/chat?username=testuser"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()
}
