package middlewares

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UserId    string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("JWT_SECRET")

func ValidateToken(signedToken string) (claim *SignedDetails, msg string) {
	// fmt.Println("token: ", signedToken)
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is not valid")
		msg = err.Error()
		return
	}

	// if claims.ExpiresAt < time.Now().Local().Unix() {
	// 	msg = fmt.Sprintf("token has expired")
	// 	// msg = err.Error()
	// 	return
	// }

	// fmt.Println(claims)
	// fmt.Println(msg)

	return claims, msg
}
