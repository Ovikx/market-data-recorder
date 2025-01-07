package jwtgen

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/go-jose/go-jose.v2"
	"gopkg.in/go-jose/go-jose.v2/jwt"
)

type apiKeyClaims struct {
	*jwt.Claims
	URI string `json:"uri"`
}

// Generates a JWT for Coinbase API requests. The URI should follow this format: `fmt.Sprintf("%s %s%s", requestMethod, requestHost, requestPath)`.
// The request host must not contain "https://".
func CoinbaseJWT(uri string) (string, error) {
	// Load the .env file variables
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %v", err)
	}

	keyName := os.Getenv("COINBASE_TRADING_API_KEY_NAME")
	privateKey := os.Getenv("COINBASE_TRADING_PRIVATE_KEY")

	// Decode private key
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", fmt.Errorf("unable to decode private key: %v", privateKey)
	}

	// Parse private key
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create new signer
	sig, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.ES256, Key: key},
		(&jose.SignerOptions{NonceSource: nonceSource{}}).WithType("JWT").WithHeader("kid", keyName),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create signer")
	}

	// Set claims
	cl := &apiKeyClaims{
		Claims: &jwt.Claims{
			Subject:   keyName,
			Issuer:    "cdp",
			NotBefore: jwt.NewNumericDate(time.Now()),
			Expiry:    jwt.NewNumericDate(time.Now().Add(2 * time.Minute)),
		},
		URI: uri,
	}

	// Sign JWT
	jwtString, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}
	return jwtString, nil

}

var max = big.NewInt(math.MaxInt64)

type nonceSource struct{}

func (n nonceSource) Nonce() (string, error) {
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return r.String(), nil
}
