package utils

import (
	"time"

	"errors"

	"github.com/TechXTT/employees-auth-svc/pkg/models"
	"github.com/golang-jwt/jwt"
)

type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

type jwtClaims struct {
	jwt.StandardClaims
	Id    string
	Email string
}

func (j *JwtWrapper) GenerateToken(admin models.Admin) (string, error) {
	claims := &jwtClaims{
		Id:    admin.ID.String(),
		Email: admin.Email,
		StandardClaims: jwt.StandardClaims{
			Issuer:    j.Issuer,
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JwtWrapper) ValidateToken(signedToken string) (*jwtClaims, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
