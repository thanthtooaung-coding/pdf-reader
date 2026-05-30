package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func IssueAccessToken(secret string, expiresHours int, userID uuid.UUID, email, roleCode string) (string, time.Time, error) {
	if secret == "" {
		return "", time.Time{}, errors.New("jwt secret is not configured")
	}
	exp := time.Now().Add(time.Duration(expiresHours) * time.Hour)
	claims := jwt.MapClaims{
		"sub":   userID.String(),
		"email": email,
		"exp":   exp.Unix(),
		"iat":   time.Now().Unix(),
	}
	if roleCode != "" {
		claims["role"] = roleCode
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, exp, nil
}

func ParseAccessToken(secret, tokenStr string) (uuid.UUID, string, string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, "", "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, "", "", errors.New("invalid token claims")
	}

	sub, _ := claims["sub"].(string)
	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, "", "", errors.New("invalid token subject")
	}

	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)
	return userID, email, role, nil
}
