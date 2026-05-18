package handlers

import (
	"backend/internal/domain"
	"backend/internal/http/response"
	"backend/internal/repo/sqlc"
	"backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserAdmin struct {
	userService    *service.UserService
	bookingService *service.BookingService
}

func NewUserAdmin(userService *service.UserService, bookingService *service.BookingService) *UserAdmin {
	return &UserAdmin{userService: userService, bookingService: bookingService}
}

// SearchUsers searches users by email substring (admin only).
// @Summary     Search users
// @Description Searches users by email substring. If query parameter 'q' is empty, lists all active users.
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       q query string false "Email substring query"
// @Success     200 {object} response.SearchUsersResponse "List of users"
// @Failure     500 {object} response.ErrorResponse "Internal server error"
// @Router      /api/v1/admin/users [get]
func (u *UserAdmin) SearchUsers(c *gin.Context) {
	q := c.Query("q")
	var users []sqlc.User
	var err error

	if q == "" {
		users, err = u.userService.ListUsers(c.Request.Context())
	} else {
		users, err = u.userService.SearchUsersByEmail(c.Request.Context(), q)
	}

	if err != nil {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("internal", "Failed to list users", nil),
		).JSON(c)
		return
	}

	docs := make([]gin.H, 0, len(users))
	for _, usr := range users {
		docs = append(docs, userToDoc(usr))
	}

	response.NewResponseBuilder(
		response.WithStatus(http.StatusOK),
		response.WithData("users", docs),
	).JSON(c)
}

// ListSellers returns all users with the seller role (admin only).
// @Summary     List sellers
// @Description Returns a list of all active users with the 'seller' role.
// @Tags        users
// @Accept      json
// @Produce     json
// @Success     200 {object} response.ListSellersResponse "List of sellers"
// @Failure     500 {object} response.ErrorResponse "Internal server error"
// @Router      /api/v1/admin/sellers [get]
func (u *UserAdmin) ListSellers(c *gin.Context) {
	sellers, err := u.userService.ListSellers(c.Request.Context())
	if err != nil {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("internal", "Failed to list sellers", nil),
		).JSON(c)
		return
	}

	docs := make([]gin.H, 0, len(sellers))
	for _, s := range sellers {
		docs = append(docs, userToDoc(s))
	}

	response.NewResponseBuilder(
		response.WithStatus(http.StatusOK),
		response.WithData("sellers", docs),
	).JSON(c)
}

// SetUserRole sets the role for a user (admin only).
// @Summary     Set user role
// @Description Sets a new role (admin, seller, client) for the specified user.
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id path string true "User ID (UUID)"
// @Param       payload body domain.SetUserRoleRequest true "Role payload"
// @Success     200 {object} response.SetUserRoleResponse "Updated user info"
// @Failure     400 {object} response.ErrorResponse "Bad request"
// @Failure     500 {object} response.ErrorResponse "Internal server error"
// @Router      /api/v1/admin/users/{id}/role [patch]
func (u *UserAdmin) SetUserRole(c *gin.Context) {
	idStr := c.Param("id")
	var uid pgtype.UUID
	if err := uid.Scan(idStr); err != nil {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("bad_request", "Invalid user ID", nil),
		).JSON(c)
		return
	}

	var req domain.SetUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("bad_request", err.Error(), nil),
		).JSON(c)
		return
	}

	updated, err := u.userService.SetUserRole(c.Request.Context(), uid, sqlc.Role(req.Role))
	if err != nil {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("internal", "Failed to set role", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithStatus(http.StatusOK),
		response.WithData("user", userToDoc(updated)),
	).JSON(c)
}

// GetBookingForSeller allows a seller or admin to view any booking by ID.
// @Summary     Get booking for seller
// @Description Fetches any booking by its ID for administrative or seller lookup.
// @Tags        bookings
// @Accept      json
// @Produce     json
// @Param       id path int true "Booking ID"
// @Success     200 {object} response.BookingResponse "Booking details"
// @Failure     400 {object} response.ErrorResponse "Bad request"
// @Failure     404 {object} response.ErrorResponse "Booking not found"
// @Router      /api/v1/seller/booking/{id} [get]
func (u *UserAdmin) GetBookingForSeller(c *gin.Context) {
	rawID := c.Param("id")
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("bad_request", "Invalid booking ID", nil),
		).JSON(c)
		return
	}

	booking, err := u.bookingService.GetBookingByID(c.Request.Context(), id)
	if err != nil {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusNotFound),
			response.WithError("not_found", "Booking not found", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithStatus(http.StatusOK),
		response.WithData("booking", toBookingDoc(booking)),
	).JSON(c)
}

func userToDoc(u sqlc.User) gin.H {
	return gin.H{
		"id":         u.ID,
		"full_name":  u.FullName,
		"email":      u.Email,
		"phone":      u.Phone.String,
		"role":       string(u.Role),
		"is_active":  u.IsActive,
		"created_at": u.CreatedAt.Time,
	}
}
