package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/shayja/go-template-api/config"
	"github.com/shayja/go-template-api/internal/entities"
	"github.com/shayja/go-template-api/internal/usecases"
)

var privateKey = []byte(config.Config("ACCESS_TOKEN_SECRET"))

type JwtUtils struct {
    UserInteractor usecases.UserInteractor
}
func GenerateJWT(user *entities.User) (string, error) {
	tokenTTL, _ := strconv.Atoi(config.Config("TOKEN_TTL"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.Id,
		"iat": time.Now().Unix(),
		"eat": time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(privateKey)
}

func ValidateJWT(context *gin.Context) error {
	token, err := getToken(context)

	if token == nil {
		return errors.New("No token provided")
	}

	if err != nil {
		return err
	}

	_, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		return nil
	}

	return errors.New("invalid client token provided")
}


func(m *JwtUtils) CurrentUser(context *gin.Context) (*entities.User, error) {
	err := ValidateJWT(context)
	if err != nil {
		return nil, err
	}

	token, _ := getToken(context)
	claims, _ := token.Claims.(jwt.MapClaims)
	userId := string(claims["id"].(string))

	user, err := m.UserInteractor.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func getToken(context *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(context)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		
		return privateKey, nil
	})
	return token, err
}

func getTokenFromRequest(context *gin.Context) string {
	bearerToken := context.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}