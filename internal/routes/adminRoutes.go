package routes

import (
	"C2S/internal/middleware"
	"C2S/internal/types"

	"github.com/gofiber/fiber/v2"
)

type AdminRoutes struct {
	store types.UserStore
}

func NewAdminRoutes(store types.UserStore) *AdminRoutes {
	return &AdminRoutes{store: store}
}

func (r *AdminRoutes) RegisterRoutes(router fiber.Router) {
	
	router.Get("/admin/users", middleware.WithJWTAuth(r.store),middleware.IsAdmin(r.store),r.store.GetAllUsers)
	router.Delete("/admin/user/:id",middleware.WithJWTAuth(r.store), middleware.IsAdmin(r.store), r.store.DeleteUser)
	router.Put("/admin/users/:id", middleware.WithJWTAuth(r.store), middleware.IsAdmin(r.store),r.store.UpdateUser)
	router.Put("/admin/user/:id",middleware.WithJWTAuth(r.store),middleware.IsAdmin(r.store),r.store.UpdateScore)
	router.Get("/admin/user/:username",middleware.WithJWTAuth(r.store),middleware.IsAdmin(r.store),r.store.GetUserByUserNameHandler)
	router.Get("/leaderboard",r.store.GetLeaderBoardHandler)
}
