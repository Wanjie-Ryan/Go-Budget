package common

import (
	"errors"
	"fmt"
	// "log"
	"os"
	"time"

	"github.com/Wanjie-Ryan/Go-Budget/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

// the struct will contain all the properties that we want to store in the JWT

type CustomJWTClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.UserModel) (*string, *string, error) {

	// what the jwt will contain
	userClaims := CustomJWTClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // will expire after 24 hours
		},
	}
	// generating token based on the user claims

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	// signing the access token using the JWT SECRET
	signedAccessToken, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // will expire after 24 hours
		},
	})

	signedRefreshToken, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, nil, err
	}

	return &signedAccessToken, &signedRefreshToken, nil

}

// to ensure the JWT is correct
// checks and extracts information from a JWT token that was sent by a client.
// 1. Verifies if the token is valid, not expired or tampered with.
// 2. extract info from token like userID
func ParseJWT(signedAccessToken string) (*CustomJWTClaims, error) {
	// go expects the custom claims structure - with fields like ID
	// the func (t *jwt.Token) provides the secret key to verify the signature, and returns JWT_SECRET
	parsedJwtAccessToken, err := jwt.ParseWithClaims(signedAccessToken, &CustomJWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		fmt.Println("error parsing access token", err)
		return nil, err
	} else if claims, ok := parsedJwtAccessToken.Claims.(*CustomJWTClaims); ok {
		return claims, nil
	} else {
		// log.Fatal("error parsing access token")
		// log.Error("error parsing access token")
		// log.Error("error parsing access token")
		fmt.Println("error parsing access token")
		return nil, errors.New("Unknown claims Type, cannot Proceed")
	}
}
