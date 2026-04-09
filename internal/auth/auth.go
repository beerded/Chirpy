package auth

import (
	"crypto/rand"
	"errors"
	"encoding/hex"
	"fmt"
	"time"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/alexedwards/argon2id"
)

type TokenType string
const TokenTypeAccess TokenType = "chirpy-access"

type AuthStringPrefix string
const BearerPrefix AuthStringPrefix = "Bearer "

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, err
}

func CheckPasswordHash(password, hash string) (bool, error) {
	res, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return res, err
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:		string(TokenTypeAccess),
		IssuedAt:	jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt:	jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:	userID.String(),
	})
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, fmt.Errorf("invalid issuer: %v", issuer)
	}

	myUUID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return myUUID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authStrings := headers.Values("Authorization")
	if len(authStrings) == 0 {
		return "", errors.New("Authorization header missing")
	}
	for _, authString := range authStrings{
		if strings.HasPrefix(authString, string(BearerPrefix)) {
			token := strings.Trim(strings.TrimPrefix(authString, string(BearerPrefix)), " ")
			return token, nil
		}
	}
	return "", errors.New("Authorization header has no bearer or it is malformed")
}

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key)
}

