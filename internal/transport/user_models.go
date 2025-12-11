package transport

type UserRegistration_req struct {
	Username string  `json:"username" example:"dinar.hadiyanto" binding:"required,min=6,max=32"`
	Email    string  `json:"email" example:"dinar.hadiyanto@outlook.com" binding:"required,email,max=255"`
	Password string  `json:"password" example:"password123" binding:"required,min=8,max=64"`
	Phone    *string `json:"phone_number" example:"+6281234567890" binding:"omitempty,e164"`
}

type User_resp struct {
	Username      string `json:"username" example:"noctisnyx"`
	Email         string `json:"email" example:"noctisnyx@noctisnyx.com"`
	Phone         string `json:"phone_number" example:"+6281234567890"`
	IsActive      bool   `json:"is_active" example:"true"`
	EmailVerified bool   `json:"email_verified" example:"false"`
	PhoneVerified bool   `json:"phone_verified" example:"false"`
}

type UserUpdate_req struct {
	Username *string `json:"username" example:"dinar.hadiyanto"`
	Email    *string `json:"email" example:"dinar.hadiyanto@outlook.com"`
	Password *string `json:"password" example:"password123"`
	Phone    *string `json:"phone_number" example:"+6281234567890"`
	IsActive *bool   `json:"is_active" example:"true"`
}

type FindUserQueryParam struct {
	ID            *string `form:"id"`
	Email         *string `form:"email"`
	Username      *string `form:"username"`
	EmailVerified *bool   `form:"email_verified"`
	PhoneVerified *bool   `form:"phone_verified"`
}
