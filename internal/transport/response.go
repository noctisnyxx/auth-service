package transport

type RestResponse struct {
	Message string      `json:"message" example:"Success or failed message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewRestResponse(msg string, data interface{}) RestResponse {
	return RestResponse{
		Message: msg,
		Data:    data,
	}
}
