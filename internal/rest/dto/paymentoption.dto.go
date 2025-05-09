package dto

type PaymentOptionAddReq struct {
	Provider string `json:"provider" form:"provider" validate:"required"`
	Account  string `json:"account" form:"account" validate:"required"`
	Name     string `json:"name" form:"name" validate:"required"`
	AsQR     bool   `json:"as_qr" form:"as_qr" validate:"required"`
}

type PaymentOptionUpdateReq struct {
	Provider string `json:"provider" form:"provider"`
	Account  string `json:"account" form:"account"`
	Name     string `json:"name" form:"name"`
	AsQR     bool   `json:"as_qr" form:"as_qr"`
}
