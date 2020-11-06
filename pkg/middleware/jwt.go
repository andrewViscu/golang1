package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GetToken(r *http.Request) (jwt.MapClaims, error) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		return nil, fmt.Errorf("Malformed token")
	}
	reqToken = splitToken[1]
	return GetAuthenticatedUser(reqToken)
}

func GetAuthenticatedUser(reqToken string) (jwt.MapClaims, error) {
	hmacSampleSecret := []byte(os.Getenv("ACCESS_SECRET"))
	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSampleSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error parsing token: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Token not valid or couldn't parse it")
	}

	return claims, nil
}
func Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := GetToken(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			next.ServeHTTP(w, r)
		}

	})
}

func CreateToken(userId string) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
