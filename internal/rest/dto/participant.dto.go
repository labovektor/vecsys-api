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
