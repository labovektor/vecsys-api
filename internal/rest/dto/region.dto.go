package dto

type RegionCreateReq struct {
	Name          string `json:"name" form:"name" validate:"required"`
	ContactName   string `json:"contact_name" form:"contact_name" validate:"required"`
	ContactNumber string `json:"contact_number" form:"contact_number" validate:"required,phone"`
	Visible       *bool  `json:"visible" form:"visible" validate:"required"`
}

type RegionEditReq struct {
	Name          string `json:"name" form:"name"`
	ContactName   string `json:"contact_name" form:"contact_name"`
	ContactNumber string `json:"contact_number" form:"contact_number,phone"`
	Visible       *bool  `json:"visible" form:"visible"`
}
