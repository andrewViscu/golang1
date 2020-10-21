package middleware

import (
	"os"
	"fmt"
	"time"
	"github.com/dgrijalva/jwt-go"
)

func GetAuthenticatedUser(tokenString string) (jwt.MapClaims, error){
	hmacSampleSecret := []byte(os.Getenv("ACCESS_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
	
		return hmacSampleSecret, nil
	})
	if err  != nil {
		return nil, fmt.Errorf("Error parsing token: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Printf("Claims type: %T", claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Token not valid or couldn't parse it")
	}
	fmt.Println("Claims: ", claims["user_id"],  claims["authorized"])
	return claims, nil
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