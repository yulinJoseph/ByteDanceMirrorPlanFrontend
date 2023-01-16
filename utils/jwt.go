package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"log"
	"strconv"
	"strings"
	"time"
)

type MyCustomClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}

func NewToken(str string) (string, error) {
	mySigningKey := []byte("im-instance")

	// Create the Claims
	claims := MyCustomClaims{
		str,
		jwt.RegisteredClaims{
			Issuer:    "",
			Subject:   "",
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			NotBefore: nil,
			IssuedAt:  nil,
			ID:        "",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return ss, nil
}

func parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("im-instance"), nil
	})
	return token, err
}

func CheckToken(tokenString string) bool {
	token, err := parseToken(tokenString)
	if err != nil {
		return false
	}
	if _, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return true
	} else {
		return false
	}
}

const (
	id = iota
	username
)

func GetFromClaims(str string, types any) string {
	ans := strings.Split(str, "+")
	switch types {
	case id:
		return ans[0]
	case username:
		return ans[1]
	default:
		return ""
	}
}

func GetIdFromToken(tokenString string) uint {
	token, err := parseToken(tokenString)
	if err != nil {
		return 0
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		id, ok := strconv.ParseUint(GetFromClaims(claims.Foo, id), 10, 64)
		if ok != nil {
			return 0
		}
		return uint(id)
	}
	return 0
}

func GetUsernameFromToken(tokenString string) string {
	token, err := parseToken(tokenString)
	if err != nil {
		return ""
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return GetFromClaims(claims.Foo, username)
	}
	return ""
}
