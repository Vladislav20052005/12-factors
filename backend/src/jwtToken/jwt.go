package jwtToken

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var hmacSampleSecret []byte

func init() {
	hmacSampleSecret = []byte(os.Getenv("AUTH_SECRET"))
	if len(hmacSampleSecret) == 0 {
		log.Fatalln("AUTH_SECRET is not defined\n")
	}
}

func GenToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString(hmacSampleSecret)
}

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}

	id, ok := claims["id"].(string)
	if !ok {
		return "", errors.New("id claim is missing or not a string")
	}

	return id, nil
}
