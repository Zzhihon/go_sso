package dto

type UserResponse struct {
	UserId      string `json:"userId"`
	Name        string `json:"name"`
	Grade       string `json:"grade"`
	MajorClass  string `json:"majorClass"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Status      string `json:"status"`
}

type UserStateResponse struct {
	Status      string `json:"status"`
	Message     string `json:"message"`
	OnlineUsers int    `json:"onlineUsers"`
	TimeStamp   int64  `json:"timeStamp"`
}
