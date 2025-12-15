package transport

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/noctisnyxx/auth-service/configs"
	"github.com/noctisnyxx/auth-service/internal/usecase"
)

type AuthHandler struct {
	logger *slog.Logger
	authUC usecase.IAuthUsecase
}

func NewAuthHandler(logger *slog.Logger, authUC usecase.IAuthUsecase) *AuthHandler {
	return &AuthHandler{
		logger: logger,
		authUC: authUC,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	var req BasicAuth_req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewRestResponse("failed to bind json", nil))
		return
	}
	authCode, err := h.authUC.BasicLogin(ctx, req.Username, req.Password, req.ClientId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, NewRestResponse("failed to login", nil))
		return
	}
	url := fmt.Sprintf("http://%s:%s/authorize?code=%s", configs.BE_HOST, configs.BE_PORT, authCode)
	c.Redirect(http.StatusFound, url)
}

func (h *AuthHandler) ExchangeCode(c *gin.Context) {

}
