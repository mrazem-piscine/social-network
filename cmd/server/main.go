package main

import (
	"log"
	"net/http"
	"strconv"

	"social-network/internal/config"
	"social-network/internal/handlers"
	"social-network/internal/middlewares"
	websockets "social-network/internal/websocket"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func main() {
	db := config.GetDB()
	defer config.CloseDB()

	middlewares.Init()

	chatManager := websockets.NewWebSocketManager()
	notificationManager := websockets.NewWebSocketNotificationManager() // ✅ Declared correctly

	r := mux.NewRouter()

	// Public Routes (No Authentication Required)
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutUser).Methods("POST")

	// Serve uploaded images
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	// WebSocket Routes
	r.HandleFunc("/ws/notifications", handlers.WebSocketNotificationHandler)
	r.HandleFunc("/ws/chat", handlers.WebSocketChatHandler)
	r.HandleFunc("/ws/group-chat", handlers.WebSocketGroupChatHandler)

	// Protected Routes (Require Authentication)
	authRoutes := r.PathPrefix("/api").Subrouter()
	authRoutes.Use(middlewares.Authenticate)

	// ✅ Post & Comment Routes
	authRoutes.HandleFunc("/comments", handlers.CreateCommentHandler).Methods("POST")
	authRoutes.HandleFunc("/comments", handlers.GetCommentsForPostHandler).Methods("GET")
	authRoutes.HandleFunc("/comments/edit", handlers.EditCommentHandler).Methods("PUT")
	authRoutes.HandleFunc("/comments", handlers.DeleteCommentHandler).Methods("DELETE")

	authRoutes.HandleFunc("/user-posts", handlers.GetUserPostsHandler).Methods("GET")
	authRoutes.HandleFunc("/posts", handlers.CreatePostHandler).Methods("POST")
	authRoutes.HandleFunc("/posts/edit", handlers.EditPostHandler).Methods("PUT")
	authRoutes.HandleFunc("/posts", handlers.DeletePostHandler).Methods("DELETE")
	// Serve uploaded images

	// ✅ Like System
	authRoutes.HandleFunc("/likes", handlers.ToggleLikeHandler).Methods("POST")
	authRoutes.HandleFunc("/likes/count", handlers.GetLikeCountHandler).Methods("GET")

	// ✅ Group Management
	authRoutes.HandleFunc("/groups", handlers.CreateGroupHandler).Methods("POST")
	authRoutes.HandleFunc("/groups/members", handlers.GetGroupMembersHandler).Methods("GET")
	authRoutes.HandleFunc("/groups/posts", handlers.CreateGroupPostHandler).Methods("POST")
	authRoutes.HandleFunc("/groups/events", handlers.CreateGroupEventHandler).Methods("POST")
	authRoutes.HandleFunc("/groups/events/rsvp", handlers.RSVPEventHandler).Methods("POST")
	authRoutes.HandleFunc("/groups/events/rsvp/count", handlers.GetRSVPCountHandler).Methods("GET")
	authRoutes.HandleFunc("/groups/events/create", handlers.CreateGroupEventHandler).Methods("POST")

	// ✅ Group Membership
	authRoutes.HandleFunc("/groups/join", handlers.RequestToJoinGroupHandler).Methods("POST") // ✅ Use only one join method
	authRoutes.HandleFunc("/groups/approve", handlers.ApproveMembershipHandler).Methods("POST")
	authRoutes.HandleFunc("/groups/reject", handlers.RejectMembershipHandler).Methods("POST")
	authRoutes.HandleFunc("/groups/leave", handlers.LeaveGroupHandler).Methods("POST")

	// ✅ Notifications
	authRoutes.HandleFunc("/notifications", handlers.GetNotificationsHandler).Methods("GET")

	authRoutes.HandleFunc("/chat/send", handlers.SendMessageHandler).Methods("POST")
	authRoutes.HandleFunc("/chat/history", handlers.GetChatHistoryHandler).Methods("GET")

	// WebSocket Chat
	r.HandleFunc("/ws/chat", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade WebSocket", http.StatusInternalServerError)
			return
		}
		chatManager.HandleConnection(conn, username, db)
	})

	// WebSocket Notifications (Fixed & Improved)
	r.HandleFunc("/ws/notifications", func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade WebSocket", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		// ✅ Store WebSocket connection for notifications
		notificationManager.Clients[userID] = &websockets.WebSocketConn{Conn: conn}

		// Send a test notification to the user
		notification := websockets.Notification{
			Type:    "notification",
			UserID:  userID,
			Message: "Welcome to the notification system!",
		}

		// ✅ Send notification through WebSocket
		err = conn.WriteJSON(notification)
		if err != nil {
			http.Error(w, "Failed to send notification", http.StatusInternalServerError)
			return
		}
	})

	log.Println("✅ Server running on :8080")
	http.ListenAndServe(":8080", r)
}
