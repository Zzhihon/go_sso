package domain

type AuthRepository interface {
	FindBy(userID string) (*Login, error)
	GenerateRefreshToken(token AuthToken) (string, error)
	RefreshTokenExists(refreshToken string) error
}
