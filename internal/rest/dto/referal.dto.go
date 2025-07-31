package dto

type ReferalCreateReq struct {
	Code          string `json:"code" validate:"required"`
	Desc          string `json:"desc" validate:"required"`
	SeatAvailable *int   `json:"seat_available" validate:"required"`
	IsDiscount    *bool  `json:"is_discount" validate:"required"`
	Discount      int    `json:"discount"`
}
