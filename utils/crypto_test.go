package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

func BenchmarkJWT(b *testing.B) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("failed to generate private key: %s\n", err)
		return
	}

	var payload []byte
	// Create signed payload
	token := jwt.New()
	token.Set(`foo`, `bar`)
	buf, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		panic(err)
	}
	b.Log("JWT BUFFER", string(buf))
	payload, err = token.Sign(jwa.RS256, privKey)
	if err != nil {
		fmt.Printf("failed to generate signed payload: %s\n", err)
		return
	}

	b.Log(string(payload))
	// Parse signed payload
	// Use jwt.ParseVerify if you want to make absolutely sure that you
	// are going to verify the signatures every time
	token, err = jwt.Parse(bytes.NewReader(payload), jwt.WithVerify(jwa.RS256, &privKey.PublicKey))
	if err != nil {
		fmt.Printf("failed to parse JWT token: %s\n", err)
		return
	}
	buf, err = json.MarshalIndent(token, "", "  ")
	if err != nil {
		fmt.Printf("failed to generate JSON: %s\n", err)
		return
	}
	fmt.Printf("%s\n", buf)

}
