package jwt

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type (
	JwtManager interface {
		Generate(userID string, tokenDuration time.Duration) (string, error)
		Verify(token string) (*jwt.Token, error)
	}

	jwtManager struct {
		secretKey string
	}

	UserClaims struct {
		jwt.StandardClaims
		UserID string `json:"UserID"`
	}
)

func NewJwtManager(secret string) JwtManager {
	return &jwtManager{
		secretKey: secret,
	}
}

func (manager *jwtManager) Generate(userID string, tokenDuration time.Duration) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenDuration).Unix(),
		},
		UserID: userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}
func (manager *jwtManager) Verify(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error : %v", t.Header["alg"])
		}
		return []byte(manager.secretKey), nil
	})
}

func GetTokenPayload(token string) (*UserClaims, error) {
	parse, _, err := new(jwt.Parser).ParseUnverified(token, &UserClaims{})
	if err != nil {
		return nil, err
	}
	payloads, ok := parse.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("failed to parse id token")
	}

	return payloads, nil
}

func GenerateForgotPasswordToken(userID string) string {
	token, err := NewJwtManager(os.Getenv("JWT_FORGOT_PASSWORD_SECRET_KEY")).Generate(userID, 1*time.Hour)
	if err != nil {
		log.Printf("Failed to generate forgot password token : %s", err)
	}
	return token
}

func GenerateAccessToken(userID string) string {
	token, err := NewJwtManager(os.Getenv("JWT_ACCESS_SECRET_KEY")).Generate(userID, 1*time.Hour)
	if err != nil {
		log.Printf("Failed to generate access token : %s", err)
	}
	return token
}

func GenerateRefreshToken(userID string) string {
	token, err := NewJwtManager(os.Getenv("JWT_REFRESH_SECRET_KEY")).Generate(userID, 730*time.Hour)
	if err != nil {
		log.Printf("Failed to generate access token : %s", err)
	}
	return token
}
