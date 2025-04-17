package dto

type NewCategoryReq struct {
	Name    string `json:"name" form:"name" validate:"required,min=2"`
	IsGroup bool   `json:"is_group" form:"is_group" validate:"required"`
	Visible bool   `json:"visible" form:"visible" validate:"required"`
}

type UpdateCategoryReq struct {
	Name    string `json:"name" form:"name"`
	IsGroup bool   `json:"is_group" form:"is_group"`
	Visible bool   `json:"visible" form:"visible"`
}
