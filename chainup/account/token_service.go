package account

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// NilToken is an empty Token used for returning and checking nil values.
var NilToken Token

// Token is a way for the Account to authenticate itself
// without having to pass Account credentials with each request.
type Token string

// NewToken returns a new Token instance.
func NewToken(token string) Token {
	return Token(token)
}

// TokenService handles generation and parsing of Tokens.
type TokenService struct {
	secret string
}

// NewTokenService returns a new TokenService instance.
func NewTokenService(secret string) *TokenService {
	return &TokenService{secret: secret}
}

// ConfigureTokenService configures the TokenService based on configuration parameters.
func ConfigureTokenService(cfg *JWTConfig) *TokenService {
	return NewTokenService(cfg.Secret)
}

// GenerateToken returns a Token with the encoded ID.
func (ts *TokenService) GenerateToken(id ID) (Token, error) {
	claims := jwt.MapClaims{
		"account_id": id.String(),
		"valid_to":   time.Now().AddDate(0, 1, 0).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(ts.secret))
	if err != nil {
		return NilToken, fmt.Errorf("sign user JWT token: %s", err)
	}

	return NewToken(tokenString), nil
}

// ParseToken returns the encoded ID from a provided Token.
func (ts *TokenService) ParseToken(token Token) (ID, error) {
	jwtToken, err := jwt.Parse(string(token), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(ts.secret), nil
	})
	if err != nil {
		return NilID, fmt.Errorf("parse user JWT token: %s", err)
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return NilID, fmt.Errorf("parse user JWT claims: %s", err)
	}

	return IDFromString(claims["account_id"].(string)), nil
}
