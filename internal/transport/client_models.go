package transport

type ClientRegistration_req struct {
	Name           string `json:"name" example:"client-app-1" binding:"required,min=6,max=64"`
	ProtocolTypeID string `json:"protocol_type_id" example:"openid-connect" binding:"required"`
}

type Client_resp struct {
	ID             string `json:"id" example:"client-app-1"`
	Name           string `json:"name" example:"client-app-1"`
	ProtocolTypeID string `json:"protocol_type_id" example:"openid-connect"`
}

type ClientDetails_resp struct {
	Client_resp
	Secret string `json:"secret" example:"uLuGusLv7qgMpAN_zAhtmX7H4W571R9bG7rmIQ60Rto="`
}

type ClientUpdate_req struct {
	Name           *string `json:"name" example:"client-app-updated"`
	ProtocolTypeID *string `json:"protocol_type_id" example:"openid-connect"`
}

type FindClientQueryParam struct {
	ProtocolTypeID *string `form:"protocol_type_id"`
}
