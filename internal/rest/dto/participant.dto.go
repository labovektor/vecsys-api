package dto

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

type RequestResetPasswordReq struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

type ResetPasswordReq struct {
	Token    string `json:"token" form:"token" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}
