package dto

type NewUpdateRequest struct {
	UserID         string `json:"userID"`
	OriginPassword string `json:"originPassword"`
	NewPassword    string `json:"newPassword"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phoneNumber"`
	Status         string `json:"status"`
	Impl           string `json:"impl"`
}
