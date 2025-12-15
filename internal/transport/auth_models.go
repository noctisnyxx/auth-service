package transport

type BasicAuth_req struct {
	Username     string
	Password     string
	ClientId     string
	ClientSecret string
}

type TokenExchange_req struct {
	Code         string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type Token_resp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
