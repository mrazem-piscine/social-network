package main

import (
	"log"
	"net/http"
	"strconv"

	"social-network/internal/config"
	"social-network/internal/handlers"
	"social-network/internal/middlewares"
	"social-network/internal/repositories"
	websockets "social-network/internal/websocket"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader            = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	chatManager         = websockets.NewChatManager()
	groupChatManager    = websockets.NewGroupChatManager()
	notificationManager = websockets.NewWebSocketNotificationManager()
)

func main() {
	// ✅ Initialize the database
	db := config.GetDB()
	defer config.CloseDB()

	// ✅ Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	groupRepo := repositories.NewGroupRepository(db)
	chatRepo := repositories.NewChatRepository(db)

	// ✅ Initialize handlers with repositories
	handlers.InitHandlers(userRepo, groupRepo, chatRepo)

	// ✅ Initialize session middleware
	middlewares.Init()

	// ✅ Create Router
	r := mux.NewRouter()

	// ✅ Public Routes (No Authentication Required)
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutUser).Methods("POST")

	// ✅ Serve uploaded images
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// ✅ WebSocket Routes (Chat & Notifications)
	r.HandleFunc("/ws/chat", WebSocketChatHandler)
	r.HandleFunc("/ws/group-chat", WebSocketGroupChatHandler)
	r.HandleFunc("/ws/notifications", WebSocketNotificationHandler)

	// ✅ Protected Routes (Require Authentication)
	authRoutes := r.PathPrefix("/api").Subrouter()
	authRoutes.Use(middlewares.Authenticate)

	// ✅ Post & Comment Routes
	authRoutes.HandleFunc("/comments", handlers.CreateCommentHandler).Methods("POST")
	authRoutes.HandleFunc("/comments", handlers.GetCommentsForPostHandler).Methods("GET")
	authRoutes.HandleFunc("/comments/edit", handlers.EditCommentHandler).Methods("PUT")
	authRoutes.HandleFunc("/comments", handlers.DeleteCommentHandler).Methods("DELETE")

	authRoutes.HandleFunc("/all-posts", handlers.GetAllPostsHandler).Methods("GET")
	authRoutes.HandleFunc("/user-posts", handlers.GetUserPostsHandler).Methods("GET")
	authRoutes.HandleFunc("/posts", handlers.CreatePostHandler).Methods("POST")
	authRoutes.HandleFunc("/posts/edit", handlers.EditPostHandler).Methods("PUT")
	authRoutes.HandleFunc("/posts", handlers.DeletePostHandler).Methods("DELETE")

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
	authRoutes.HandleFunc("/groups/join", handlers.RequestToJoinGroupHandler).Methods("POST")
	authRoutes.HandleFunc("/groups/approve", handlers.ApproveMembershipHandler).Methods("POST")
	authRoutes.HandleFunc("/groups/reject", handlers.RejectMembershipHandler).Methods("POST")
	authRoutes.HandleFunc("/groups/leave", handlers.LeaveGroupHandler).Methods("POST")
	authRoutes.HandleFunc("/groups/chat/history", handlers.GetGroupChatHistoryHandler).Methods("GET")
	authRoutes.HandleFunc("/groups/chat/send", handlers.SendGroupChatMessageHandler).Methods("POST")

	// ✅ Notifications
	authRoutes.HandleFunc("/notifications", handlers.GetNotificationsHandler).Methods("GET")
	authRoutes.HandleFunc("/notifications", handlers.GetNotificationsHandler).Methods("GET")
	authRoutes.HandleFunc("/notifications/send", handlers.SendNotificationHandler).Methods("POST") // <== Add this!
	// authRoutes.HandleFunc("/notifications/read", handlers.MarkNotificationsAsReadHandler).Methods("PUT")

	// ✅ Private Chat
	authRoutes.HandleFunc("/chat/send", handlers.SendMessageHandler).Methods("POST")
	authRoutes.HandleFunc("/chat/history", handlers.GetChatHistoryHandler).Methods("GET")

	log.Println("✅ Server running on :8080")
	http.ListenAndServe(":8080", r)
}

// ✅ WebSocket Handler for Direct Chat
func WebSocketChatHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade WebSocket", http.StatusInternalServerError)
		return
	}

	chatManager.HandleChatConnection(conn, userID)
}

// ✅ WebSocket Handler for Group Chat
func WebSocketGroupChatHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := strconv.Atoi(r.URL.Query().Get("group_id"))
	if err != nil || groupID == 0 {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade WebSocket", http.StatusInternalServerError)
		return
	}

	groupChatManager.JoinGroupChat(conn, groupID, userID)
}

// ✅ WebSocket Handler for Notifications
func WebSocketNotificationHandler(w http.ResponseWriter, r *http.Request) {
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
	notification := websockets.Notification{
		Type:    "notification",
		UserID:  userID,
		Message: "Welcome to the notification system!",
	}

	err = conn.WriteJSON(notification)
	if err != nil {
		http.Error(w, "Failed to send notification", http.StatusInternalServerError)
	}
}
