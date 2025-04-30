package dto

import "github.com/labovector/vecsys-api/entity"

type ParticipantLoginReq struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type ParticipantSignUpReq struct {
	EventId  string `json:"event_id" form:"event_id" validate:"required"`
	Name     string `json:"name" form:"name" validate:"required"`
	Email    string ` json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type ForgotPasswordReq struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

type ResetPasswordReq struct {
	Token    string `json:"token" form:"token" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type PickCategoryAndRegionReq struct {
	CategoryId string `json:"category_id" form:"category_id" validate:"required"`
	RegionId   string `json:"region_id" form:"region_id" validate:"required"`
}

type ClaimReferalReq struct {
	Code string `json:"code" form:"code" validate:"required"`
}

type PaymentReq struct {
	PaymentOptionId string `json:"payment_option_id" form:"payment_option_id" validate:"required"`
	AccountNumber   string `json:"account_number" form:"account_number" validate:"required"`
	AccountName     string `json:"account_name" form:"account_name" validate:"required"`
	TransferDate    string `json:"transfer_date" form:"transfer_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

type AddInstitutionReq struct {
	Name            string `json:"name" form:"name" validate:"required"`
	Email           string `json:"email" form:"email" validate:"required,email"`
	PendampingName  string `json:"pendamping_name" form:"phone" validate:"required"`
	PendampingPhone string `json:"pendamping_phone" form:"pendamping_phone" validate:"required,phone"`
}

type PickInstitutionReq struct {
	InstitutionId string `json:"institution_id" form:"institution_id" validate:"required"`
}

type AddMemberReq struct {
	Name     string        `json:"name" form:"name" validate:"required"`
	Gender   entity.Gender `json:"gender" form:"gender" validate:"required"`
	Email    string        `json:"email" form:"email" validate:"required,email"`
	Phone    string        `json:"phone" form:"phone" validate:"required,phone"`
	IDNumber string        `json:"id_number" form:"id_number" validate:"required"`
}
