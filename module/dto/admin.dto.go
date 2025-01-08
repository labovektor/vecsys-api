package dto

type AdminLoginReq struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type AdminSignUpReq struct {
	Username    string `json:"username" form:"username"`
	DisplayName string `json:"display_name" form:"display_name"`
	Email       string `json:"email" form:"email"`
	Password    string `json:"password" form:"password"`
}

type AdminEditReq struct {
	DisplayName string `json:"display_name" form:"display_name"`
	Email       string `json:"email" form:"email"`
}
