package jwtAuth

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"serverFordownDrive/config"
	"time"
)

func GeenerateJWT(user_id string) (string, error) {
	mySigningKey := []byte(config.GetJWTSecret())
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		println(err.Error())
		return "", err
	}
	return tokenString, nil
}

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		authCookie, err := request.Cookie("token")
		if err != nil {
			println(err.Error())
			return
		}
		mySigningKey := []byte(config.GetJWTSecret())

		token, err := jwt.Parse(authCookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was error parsing token")
			}
			return mySigningKey, nil
		})

		if err != nil {
			println("your Token has been expired")
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId := fmt.Sprintf("%s", claims["user"])
			request.Header.Set("user", userId)
			handler.ServeHTTP(writer, request)
			return
		}
	}
}

func IsAuthorized2(request *http.Request) (string, error) {

	authCookie, err := request.Cookie("token")
	if err != nil {
		println(err.Error())
		return "", err
	}
	mySigningKey := []byte(config.GetJWTSecret())

	token, err := jwt.Parse(authCookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was error parsing token")
		}
		return mySigningKey, nil
	})

	if err != nil {
		println("your Token has been expired")
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := fmt.Sprintf("%s", claims["user"])
		return userId, nil
	}

	return "", fmt.Errorf("passes all ifs in isautherized2")

}
