package authentication

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type JWTService struct {
	signKey []byte
	issuer  string
	expire  time.Duration
}

func NewJWTService(secret, issuer string, expire time.Duration) *JWTService {
	return &JWTService{
		signKey: []byte(secret),
		issuer:  issuer,
		expire:  expire,
	}
}

func (s *JWTService) CreateToken(body *User) (string, error) {
	logrus.Infof("User: %+v", body)
	now := time.Now()
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		CustomClaims{
			User: body,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    s.issuer,
				ExpiresAt: jwt.NewNumericDate(now.Add(s.expire)),
				NotBefore: jwt.NewNumericDate(now),
				ID:        uuid.NewString(),
			},
		},
	)
	return token.SignedString(s.signKey)
}

func (s *JWTService) ParseToken(tokenString string) (*User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.signKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims.User, nil
}

func (s *JWTService) GetNeedRefreshToken(tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.signKey, nil
	})
	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return false, errors.New("invalid token")
	}

	// 如果小于5天则返回true
	if time.Now().Add(5 * time.Hour * 24).After(claims.ExpiresAt.Time) {
		return true, nil
	} else {
		return false, nil
	}
}

func (s *JWTService) RefreshToken(tokenString string) (string, error) {
	body, err := s.ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	return s.CreateToken(body)
}

func (s *JWTService) Initialize() error {
	return nil
}

type CustomClaims struct {
	User *User
	jwt.RegisteredClaims
}
