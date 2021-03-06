package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	ENCRYPTION_SALT         = "jh9J6nGvRyFznCjHJXgaLM"
	PASSWORD_SALT           = "ByBDCG2sAYK1IMP"
	JWT_SIGNATURE_ALGORITHM = jwa.HS256
	PURP_ACCESS_TOKEN       = "PURP_ACCESS_TOKEN"
	PURP_FORGOT_PASSWORD    = "PURP_FORGOT_PASSWORD"
)

func GeneratePasswordHash(email, password string) (string, error) {
	combined := email + ":" + password + ":" + PASSWORD_SALT
	bytes, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	return string(bytes), err
}

func ValidatePasswordHash(incoming, validator string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(validator), []byte(incoming))
	return err == nil
}

func GenerateWebToken(id, token, purpose string) ([]byte, error) {
	// log := loggerFactory.CreateLog().SetFunctionName("GenerateWebToken").SetStartTime()
	// defer log.SetFinishTime().WriteAndDeleteLog()

	tokenJwt := jwt.New()
	tokenJwt.Set(`ID`, id)
	tokenJwt.Set(`InternalToken`, token)
	tokenJwt.Set(`Purpose`, purpose)
	payload, err := tokenJwt.Sign(JWT_SIGNATURE_ALGORITHM, []byte(ENCRYPTION_SALT))
	if err != nil {
		return []byte{}, err
	}
	return payload, nil
}
func ValidateWebToken(webToken []byte) bool {
	options := jwt.WithVerify(JWT_SIGNATURE_ALGORITHM, []byte(ENCRYPTION_SALT))
	_, err := jwt.Parse(bytes.NewReader(webToken), options)
	if err != nil {
		fmt.Printf("failed to parse JWT token: %s\n", err)
		return false
	}
	return true
}

// Get webtoken from a http request, will return with userprofile and error
func GetWebToken(r *http.Request) (User, error) {
	var err error
	auth, ok := r.Header["Authorization"]
	if !ok {
		err = errors.New("ERR_CANNOT_PARSE_HEADER")
		return User{}, err
	}

	splitAuth := strings.Split(auth[0], " ")
	if splitAuth[0] != "Bearer" {
		err = errors.New("ERR_WRONG_AUTHORIZATION")
		return User{}, err
	}
	return VerifyToken(splitAuth[1], PURP_ACCESS_TOKEN)
}

func VerifyToken(incomingToken, tokenPurpose string) (User, error) {
	convertedToken, err := base64.StdEncoding.DecodeString(incomingToken)
	if err != nil {
		err = errors.New("ERR_WRONG_AUTHORIZATION")
		return User{}, err
	}

	options := jwt.WithVerify(JWT_SIGNATURE_ALGORITHM, []byte(ENCRYPTION_SALT))
	token, err := jwt.Parse(bytes.NewReader(convertedToken), options)
	if err != nil {
		err = errors.New("ERR_WRONG_AUTHORIZATION")
		return User{}, err
	}
	userId, _ := token.Get("ID")
	userToken, _ := token.Get("InternalToken")

	var user User
	user.GetUser(userId.(string))

	var storedToken string
	switch tokenPurpose {
	case PURP_ACCESS_TOKEN:
		storedToken = *user.AccessToken
	case PURP_FORGOT_PASSWORD:
		storedToken = *user.PassToken
	}

	if storedToken != userToken {
		err = errors.New("ERR_WRONG_AUTHORIZATION")
		return User{}, err
	}
	return user, nil
}
