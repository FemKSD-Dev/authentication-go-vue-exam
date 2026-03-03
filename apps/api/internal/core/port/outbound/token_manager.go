package outbound

type TokenClaims struct {
	UserID   string
	Username string
}

type TokenManager interface {
	Issue(userID, username string) (string, string, error)
	Verify(token string) (*TokenClaims, error)
}
