package plugins

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	ErrTokenExpired = errors.New("token is expired")
)

type Tokens struct {
	AccessToken  TokensData
	RefreshToken TokensData
}

type TokensData struct {
	Token  string
	Secret []byte
	TTL    int
}

func (t *Tokens) GenSecret() error {
	t.AccessToken.Secret = make([]byte, 32)
	t.RefreshToken.Secret = make([]byte, 32)
	if _, err := rand.Read(t.AccessToken.Secret); err != nil {
		return err
	}
	if _, err := rand.Read(t.RefreshToken.Secret); err != nil {
		return err
	}
	if t.AccessToken.TTL == 0 {
		t.AccessToken.TTL = 10
	}
	if t.RefreshToken.TTL == 0 {
		t.RefreshToken.TTL = 60
	}
	return nil
}
func (t *Tokens) GetToken() error {
	var err error
	t.AccessToken.Token, err = createToken(t.AccessToken.Secret, t.AccessToken.TTL)
	if err != nil {
		return fmt.Errorf("error creating JWT token: %w", err)
	}
	t.RefreshToken.Token, err = createToken(t.RefreshToken.Secret, t.RefreshToken.TTL)
	if err != nil {
		return fmt.Errorf("error creating refresh token: %w", err)
	}
	return nil
}

func (t *Tokens) ValidateToken(tokenData TokensData) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenData.Token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("no valid signing method")
		}
		return tokenData.Secret, nil
	})
	if err != nil {
		if err.Error() == "Token is expired" {
			uuid, err := setClaims(token)
			if err != nil {
				return uuid, fmt.Errorf("failed to extract claims from expired token: %w", err)
			}
			return uuid, ErrTokenExpired
		}
		return uuid.Nil, fmt.Errorf("failed to parse token: %w", err)
	}
	return setClaims(token)
}

func createToken(secret []byte, tokenTTL int) (string, error) {
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":  "narria",
			"iss":  "diary",
			"role": "system",
			"exp":  time.Now().Add(time.Duration(tokenTTL * int(time.Minute))).Unix(),
			"iat":  time.Now().Unix(),
		}).SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}
	return tokenString, nil
}

func setClaims(token *jwt.Token) (uuid.UUID, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	id, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("cant parse id from jwt token")
	}
	return uuid.MustParse(id), nil
}
