package utils

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/models"
)

const jwtSecret = "a;W>)@$p9f3Ugn!"

type Claims struct {
	*jwt.StandardClaims
	models.User
}

func getExpiry() time.Time {
	return time.Now().Add(time.Hour * 24)
}

func DecodeToken(token string) (*models.User, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims := parsed.Claims.(*Claims)

	return &claims.User, nil
}

func CreateToken(user models.User) (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims = &Claims{
		&jwt.StandardClaims{
			ExpiresAt: getExpiry().Unix(),
		},
		user,
	}

	return t.SignedString([]byte(jwtSecret))
}

func CreateJWTCookie(user models.User) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "authtoken"
	token, err := CreateToken(user)
	if err != nil {
		return nil
	}

	cookie.Value = token
	cookie.Expires = getExpiry()

	return cookie
}

func GetJWTCookie(cookies []*http.Cookie) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == "authtoken" {
			return cookie
		}
	}

	return nil
}
