package dto

type LoginRequest struct {
	UserID   string `json:"userID"`
	Password string `json:"password"`
}

type VerifyRequest struct {
	Token string `json:"token"`
}

type RefreshRequest struct {
	AccessToken  string `json:"accessToken"`
	RefrestToken string `json:"refrestToken"`
}
