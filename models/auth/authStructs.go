package auth

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt"
)

type AuthRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}
type AuthResponse struct {
	Error  string `json:"error"`
	Token  string `json:"token"`
	Client int    `json:"client"`
	User   string `json:"user"`
}

type Claim struct {
	Credentials *AuthRequest
	jwt.StandardClaims
}

var (
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
)
