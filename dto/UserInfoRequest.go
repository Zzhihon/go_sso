package dto

type NewUpdateRequest struct {
	UserID      string `json:"userID"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	Impl        string `json:"impl"`
	Code        string `json:"code"`
}

type CheckEmailRequest struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
}

type EmailData struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Time string `json:"time"`
}

type GetAllUsers struct {
	Status   string `json:"status"`
	Page     int    `json:"page"`
	PageSize int    `json:"Size"`
}

type OnlineUsers struct {
	UserID string `json:"userID"`
}
