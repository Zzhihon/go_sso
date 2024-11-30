package dto

type NewUpdateRequest struct {
	UserID      string `json:"userID"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Status      string `json:"status"`
	Impl        string `json:"impl"`
}
