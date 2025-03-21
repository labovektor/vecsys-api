package dto

type APIResponse struct {
	Status *APIStatus `json:"status"`
	Data   any        `json:"data"`
}

type APIStatus struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (as APIStatus) WithMessage(message string) *APIStatus {
	as.Message = message
	return &as
}

var (
	SuccessStatus = &APIStatus{
		Success: true,
		Message: "Success",
	}
	ErrorStatus = &APIStatus{
		Success: false,
		Message: "Error",
	}
)
