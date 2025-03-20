package dto

type ParticipantLoginReq struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type ParticipantSignUpReq struct {
	EventId  string `json:"event_id" form:"event_id"`
	Name     string `json:"name" form:"name"`
	Email    string ` json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
