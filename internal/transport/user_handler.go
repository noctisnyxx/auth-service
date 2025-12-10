package transport

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/noctisnyxx/auth-service/internal/domain"
	"github.com/noctisnyxx/auth-service/internal/usecase"
)

type UserHandler struct {
	logger *slog.Logger
	userUC usecase.IUserUsecase
}

func NewUserHandler(logger *slog.Logger, userUC usecase.IUserUsecase) *UserHandler {
	return &UserHandler{
		logger: logger,
		userUC: userUC,
	}
}

// @Summary		Register a new user
// @Description	Register a new user with the provided details
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			user	body		UserRegistration_req	true	"User registration details"
// @Success		200		{object}	RestResponse			"User registered successfully"
// @Failure		400		{object}	RestResponse			"Bad request"
// @Failure		500		{object}	RestResponse			"Internal server error"
// @Router			/user/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	var req UserRegistration_req
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind json", "error", err)
		c.JSON(http.StatusBadRequest, NewRestResponse("failed to bind json", nil))
		return
	}

	var user_dom domain.User
	user_dom.Username = req.Username
	user_dom.Email = req.Email
	user_dom.Password = req.Password
	user_dom.Phone = req.Phone

	if err := h.userUC.RegisterUser(ctx, &user_dom); err != nil {
		h.logger.Error("failed to register user", "error", err)
		c.JSON(http.StatusInternalServerError, NewRestResponse("failed to register user", nil))
		return
	}

	c.JSON(http.StatusCreated, NewRestResponse("user registered successfully", nil))

}

// @Summary		Find users
// @Description	Find users by ID, email, or username
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			id				query		string							false	"User ID"
// @Param			email			query		string							false	"User email"
// @Param			username		query		string							false	"User username"
// @Param			email_verified	query		bool							false	"User email verification status"
// @Param			phone_verified	query		bool							false	"User phone verification status"
// @Success		200				{object}	RestResponse{data=[]User_resp}	"Users found successfully"
// @Failure		400				{object}	RestResponse					"Bad request"
// @Failure		500				{object}	RestResponse					"Internal server error"
// @Router			/user [get]
func (h *UserHandler) FindUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	var query FindUserQueryParam
	if err := c.ShouldBindQuery(&query); err != nil {
		h.logger.Error("failed to bind query", "error", err)
		c.JSON(http.StatusBadRequest, NewRestResponse("failed to bind query", nil))
		return
	}

	var domainQuery domain.FindUserQuery
	domainQuery.ID = query.ID
	domainQuery.Email = query.Email
	domainQuery.Username = query.Username
	domainQuery.EmailVerified = query.EmailVerified
	domainQuery.PhoneVerified = query.PhoneVerified

	users, err := h.userUC.FindUser(ctx, domainQuery)
	if err != nil {
		h.logger.Error("failed to find users", "error", err)
		c.JSON(http.StatusInternalServerError, NewRestResponse("failed to find users", nil))
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusOK, NewRestResponse("no users found", nil))
		return
	}

	var usersResp []User_resp
	for _, user := range users {
		var Phone string
		if user.Phone != nil {
			Phone = *user.Phone
		}
		usersResp = append(usersResp, User_resp{
			Username:      user.Username,
			Email:         user.Email,
			Phone:         Phone,
			IsActive:      user.IsActive,
			EmailVerified: user.EmailVerified,
			PhoneVerified: user.PhoneVerified,
		})
	}

	c.JSON(http.StatusOK, NewRestResponse("users found successfully", usersResp))
}
