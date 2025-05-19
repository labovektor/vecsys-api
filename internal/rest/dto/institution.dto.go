package dto

type InstitutionAddReq struct {
	Name            string `json:"name" form:"name" validate:"required"`
	Email           string `json:"email" form:"email" validate:"required,email"`
	PendampingName  string `json:"pendamping_name" form:"phone" validate:"required"`
	PendampingPhone string `json:"pendamping_phone" form:"pendamping_phone" validate:"required,phone"`
}

type InstitutionUpdateReq struct {
	Name            string `json:"name" form:"name"`
	Email           string `json:"email" form:"email"`
	PendampingName  string `json:"pendamping_name" form:"phone"`
	PendampingPhone string `json:"pendamping_phone" form:"pendamping_phone"`
}
