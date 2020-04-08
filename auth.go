package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/humamfauzi/go-registration/utils"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	ENCRYPTION_SALT         = "jh9J6nGvRyFznCjHJXgaLM"
	JWT_SIGNATURE_ALGORITHM = jwa.HS256
)

func GeneratePasswordHash(email, password string) (string, error) {
	combined := email + ":" + password
	bytes, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	return string(bytes), err
}

func ValidatePasswordHash(incoming, validator string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(validator), []byte(incoming))
	return err == nil
}

func GenerateWebToken(Id string) ([]byte, error) {
	log := loggerFactory.CreateLog().SetFunctionName("GenerateWebToken").SetStartTime()
	defer log.SetFinishTime().WriteAndDeleteLog()

	token := jwt.New()
	token.Set(`ID`, Id)
	token.Set(`InternalToken`, utils.GenerateUUID("token", 4))
	token.Set(`ValidUntil`, time.Now().Add(time.Hour*time.Duration(24)))
	payload, err := token.Sign(JWT_SIGNATURE_ALGORITHM, []byte(ENCRYPTION_SALT))
	if err != nil {
		return []byte{}, err
	}
	return payload, nil
}

func ValidateWebToken(webToken []byte) bool {
	options := jwt.WithVerify(JWT_SIGNATURE_ALGORITHM, []byte(ENCRYPTION_SALT))
	token, err := jwt.Parse(bytes.NewReader(webToken), options)
	if err != nil {
		fmt.Printf("failed to parse JWT token: %s\n", err)
		return false
	}
	return true
}
