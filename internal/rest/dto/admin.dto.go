package dto

type AdminLoginReq struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type AdminSignUpReq struct {
	Username    string `json:"username" form:"username" validate:"required"`
	DisplayName string `json:"display_name" form:"display_name" validate:"required,min=5,max=100"`
	Email       string `json:"email" form:"email" validate:"required,email"`
	Password    string `json:"password" form:"password" validate:"required"`
}

type AdminEditReq struct {
	DisplayName string `json:"display_name" form:"display_name" validate:"required,min=5,max=100"`
	Email       string `json:"email" form:"email" validate:"required,email"`
}
