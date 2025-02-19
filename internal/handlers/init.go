package handlers

import "social-network/internal/repositories"

var (
	userRepo  *repositories.UserRepository
	groupRepo *repositories.GroupRepository
	chatRepo  *repositories.ChatRepository
)

// InitHandlers initializes the handlers with necessary repositories
func InitHandlers(uRepo *repositories.UserRepository, gRepo *repositories.GroupRepository, cRepo *repositories.ChatRepository) {
	userRepo = uRepo
	groupRepo = gRepo
	chatRepo = cRepo
}
