package transport

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/noctisnyxx/auth-service/internal/domain"
	"github.com/noctisnyxx/auth-service/internal/usecase"
)

type ClientHandler struct {
	logger   *slog.Logger
	clientUC usecase.IClientUsecase
}

func NewClientHandler(logger *slog.Logger, clientUC usecase.IClientUsecase) *ClientHandler {
	return &ClientHandler{
		logger:   logger,
		clientUC: clientUC,
	}
}

// @Summary		Register a new client
// @Description	Register a new client with the provided details
// @Tags			Client
// @Accept			json
// @Produce		json
// @Param			client	body		ClientRegistration_req	true	"Client registration details"
// @Success		201		{object}	RestResponse			"Client registered successfully"
// @Failure		400		{object}	RestResponse			"Bad request"
// @Failure		500		{object}	RestResponse			"Internal server error"
// @Router			/clients [post]
func (h *ClientHandler) RegisterClient(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	var req ClientRegistration_req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewRestResponse("failed to bind json", nil))
		return
	}

	var client_dom domain.Client
	client_dom.Name = req.Name
	client_dom.ProtocolTypeID = req.ProtocolTypeID

	if err := h.clientUC.Register(ctx, &client_dom); err != nil {
		if errors.Is(err, domain.ErrDuplicateClientName) {
			c.JSON(http.StatusBadRequest, NewRestResponse(err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, NewRestResponse("failed to register client", nil))
		return
	}

	c.JSON(http.StatusCreated, NewRestResponse("client registered successfully", gin.H{
		"id":     client_dom.ID,
		"secret": client_dom.Secret,
	}))
}

// @Summary		Find clients
// @Description	Find clients by ID, name, or protocol type
// @Tags			Client
// @Accept			json
// @Produce		json
// @Param			protocol_type_id	query		string								false	"Client protocol type ID"
// @Success		200					{object}	RestResponse{data=[]Client_resp}	"Clients found successfully"
// @Failure		400					{object}	RestResponse						"Bad request"
// @Failure		500					{object}	RestResponse						"Internal server error"
// @Router			/clients [get]
func (h *ClientHandler) FindClients(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	var query FindClientQueryParam
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, NewRestResponse("failed to bind query", nil))
		return
	}

	var domainQuery domain.FindClientQuery
	domainQuery.ProtocolTypeID = query.ProtocolTypeID

	clients, err := h.clientUC.Find(ctx, domainQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewRestResponse("failed to find clients", nil))
		return
	}
	if len(clients) == 0 {
		c.JSON(http.StatusOK, NewRestResponse("no clients found", nil))
		return
	}

	var clientsResp []Client_resp
	for _, client := range clients {
		clientsResp = append(clientsResp, Client_resp{
			ID:             client.ID,
			Name:           client.Name,
			ProtocolTypeID: client.ProtocolTypeID,
		})
	}

	c.JSON(http.StatusOK, NewRestResponse("clients found successfully", clientsResp))
}

// @Summary		Find client details
// @Description	Find client details by ID
// @Tags			Client
// @Accept			json
// @Produce		json
// @Param			id	path		string									true	"Client ID"
// @Success		200	{object}	RestResponse{data=ClientDetails_resp}	"Client found successfully"
// @Failure		400	{object}	RestResponse							"Bad request"
// @Failure		404	{object}	RestResponse							"Client not found"
// @Failure		500	{object}	RestResponse							"Internal server error"
// @Router			/clients/{id} [get]
func (h *ClientHandler) FindClientDetails(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	id := c.Param("id")

	clients, err := h.clientUC.Find(ctx, domain.FindClientQuery{ID: &id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewRestResponse("failed to find client", nil))
		return
	}
	if len(clients) == 0 {
		c.JSON(http.StatusNotFound, NewRestResponse("client not found", nil))
		return
	}

	client := clients[0]
	clientResp := ClientDetails_resp{
		Client_resp: Client_resp{
			ID:             client.ID,
			Name:           client.Name,
			ProtocolTypeID: client.ProtocolTypeID,
		},
		Secret: client.Secret,
	}
	c.JSON(http.StatusOK, NewRestResponse("client found successfully", clientResp))
}

// @Summary		Update a client
// @Description	Update a client's details by ID
// @Tags			Client
// @Accept			json
// @Produce		json
// @Param			id		path		string				true	"Client ID"
// @Param			client	body		ClientUpdate_req	true	"Client update details"
// @Success		200		{object}	RestResponse		"Client updated successfully"
// @Failure		400		{object}	RestResponse		"Bad request"
// @Failure		404		{object}	RestResponse		"Client not found"
// @Failure		500		{object}	RestResponse		"Internal server error"
// @Router			/clients/{id} [put]
func (h *ClientHandler) UpdateClient(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	id := c.Param("id")
	var req ClientUpdate_req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewRestResponse("failed to bind json", nil))
		return
	}

	var clientUpdate_dom domain.ClientUpdate
	clientUpdate_dom.Name = req.Name
	clientUpdate_dom.ProtocolTypeID = req.ProtocolTypeID

	if err := h.clientUC.Update(ctx, id, clientUpdate_dom); err != nil {
		c.JSON(http.StatusInternalServerError, NewRestResponse("failed to update client", nil))
		return
	}

	c.JSON(http.StatusOK, NewRestResponse("client updated successfully", nil))
}

// @Summary		Delete a client
// @Description	Delete a client by ID
// @Tags			Client
// @Accept			json
// @Produce		json
// @Param			id	path		string			true	"Client ID"
// @Success		200	{object}	RestResponse	"Client deleted successfully"
// @Failure		400	{object}	RestResponse	"Bad request"
// @Failure		404	{object}	RestResponse	"Client not found"
// @Failure		500	{object}	RestResponse	"Internal server error"
// @Router			/clients/{id} [delete]
func (h *ClientHandler) DeleteClient(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	id := c.Param("id")

	if err := h.clientUC.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, NewRestResponse("failed to delete client", nil))
		return
	}

	c.JSON(http.StatusOK, NewRestResponse("client deleted successfully", nil))
}
