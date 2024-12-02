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
}
